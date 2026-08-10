# Copyright 2022 iLogtail Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.22)
project(driver)

if (LINUX)
    set(ENABLE_PROFILE OFF CACHE BOOL "Enable profile" FORCE)
    set(BUILD_TESTING OFF CACHE BOOL "Build test cases" FORCE)
    set(CMAKE_INSTALL_PREFIX ${DEFAULT_DEPS_ROOT})
    
    # Save the original C flags
    set(ORIGINAL_CMAKE_C_FLAGS ${CMAKE_C_FLAGS})
    
    # Remove -Werror for coolbpf compilation (only if ENABLE_WERROR is ON)
    if (ENABLE_WERROR)
        string(REPLACE "-Werror" "" CMAKE_C_FLAGS "${CMAKE_C_FLAGS}")
    endif()

    # LoongCollector Linux: build AgentSight (cargo → libagentsight.so + agentsight.h) with coolbpf.
    # coolbpf/CMakeLists.txt uses option(ENABLE_AGENTSIGHT OFF) on the same cache key as
    # core/CMakeLists.txt cmake_dependent_option. Without CACHE+FORCE here, a stale OFF (e.g. leftover
    # cache from a standalone coolbpf configure) keeps add_subdirectory(agentsight) disabled and
    # no libagentsight target is created.
    if (NOT DEFINED ENABLE_AGENTSIGHT)
        set(ENABLE_AGENTSIGHT OFF)
    endif ()
    set(ENABLE_AGENTSIGHT "${ENABLE_AGENTSIGHT}" CACHE BOOL "Build AgentSight library (libagentsight.so + agentsight.h)" FORCE)
    message(STATUS "ebpf/driver: ENABLE_AGENTSIGHT=${ENABLE_AGENTSIGHT} (coolbpf agentsight / libagentsight)")

    # ─── TEMPORARY: pre-upstream agentsight patch ────────────────────────────
    # Applied to the vendored coolbpf WORKING TREE before it is configured, because
    # add_subdirectory below compiles files on disk rather than a git ref. This lets us
    # build and test the fix while the upstream anolisa PR is still open, without
    # inventing a coolbpf submodule commit that the CI cannot fetch from gitee.
    #
    # REMOVE this block (and the patch file) once the fix lands upstream and the coolbpf
    # submodule pointer is bumped instead. Tracking: alibaba/anolisa
    # "fix(sight): report pids in the observer namespace".
    set(_AGENTSIGHT_PATCH "${CMAKE_SOURCE_DIR}/_thirdparty/patches/agentsight-observer-pidns.patch")
    if (ENABLE_AGENTSIGHT AND EXISTS "${_AGENTSIGHT_PATCH}")
        find_program(GIT_EXECUTABLE git REQUIRED)
        set(_COOLBPF_DIR "${CMAKE_SOURCE_DIR}/_thirdparty/coolbpf")
        # `git apply --reverse --check` succeeds only when the patch is ALREADY applied,
        # which keeps re-configures (and warm CI caches) idempotent.
        execute_process(
                COMMAND "${GIT_EXECUTABLE}" apply --reverse --check -p1 "${_AGENTSIGHT_PATCH}"
                WORKING_DIRECTORY "${_COOLBPF_DIR}"
                RESULT_VARIABLE _AGENTSIGHT_PATCH_PRESENT
                OUTPUT_QUIET ERROR_QUIET)
        if (_AGENTSIGHT_PATCH_PRESENT EQUAL 0)
            message(STATUS "ebpf/driver: agentsight observer-pidns patch already applied, skipping")
        else ()
            execute_process(
                    COMMAND "${GIT_EXECUTABLE}" apply -p1 "${_AGENTSIGHT_PATCH}"
                    WORKING_DIRECTORY "${_COOLBPF_DIR}"
                    RESULT_VARIABLE _AGENTSIGHT_PATCH_RC
                    ERROR_VARIABLE _AGENTSIGHT_PATCH_ERR)
            # Fail loudly: a silently unpatched build would look like the fix does not work.
            if (NOT _AGENTSIGHT_PATCH_RC EQUAL 0)
                message(FATAL_ERROR
                        "ebpf/driver: failed to apply ${_AGENTSIGHT_PATCH} in ${_COOLBPF_DIR}\n"
                        "  ${_AGENTSIGHT_PATCH_ERR}\n"
                        "  The coolbpf submodule may have moved past the patch base. "
                        "Regenerate the patch, or drop it if the fix is already upstream.")
            endif ()
            message(STATUS "ebpf/driver: applied agentsight observer-pidns patch to coolbpf")
        endif ()
    endif ()

    # Now add the coolbpf subdirectory (adds agentsight/ when ENABLE_AGENTSIGHT is on).
    add_subdirectory(${CMAKE_SOURCE_DIR}/_thirdparty/coolbpf ${CMAKE_BINARY_DIR}/_thirdparty/coolbpf)

    # coolbpf AgentSight: vendored openssl-sys runs nested `make` and breaks when GNU make forwards a
    # jobserver (parallel top-level make). Prefer not to patch coolbpf/src/agentsight/CMakeLists.txt;
    # USES_TERMINAL reduces jobserver handoff into the cargo/openssl make chain (CMake 3.2+).
    if (TARGET libagentsight)
        set_target_properties(libagentsight PROPERTIES USES_TERMINAL TRUE)
    endif ()

    # Restore the original C flags for the rest of the build
    set(CMAKE_C_FLAGS ${ORIGINAL_CMAKE_C_FLAGS})
    
    set(_coolbpf_install_deps coolbpf coolbpf_static)
    # coolbpf installs libagentsight.so via install(FILES ...): the producer may be UTILITY (cargo) or a real SHARED
    # target. Either way install_coolbpf must depend on it so `cmake --install` runs after the .so exists.
    if (ENABLE_AGENTSIGHT AND TARGET libagentsight)
        list(APPEND _coolbpf_install_deps libagentsight)
    elseif (ENABLE_AGENTSIGHT)
        message(STATUS "ENABLE_AGENTSIGHT=ON but no target 'libagentsight' in coolbpf; install_coolbpf skips it.")
    endif ()
    add_custom_target(
        build_libprofiler ALL
        COMMAND CARGO_TARGET_DIR=${CMAKE_BINARY_DIR}/_thirdparty/coolbpf/src/profiler cargo build -r
        WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}/_thirdparty/coolbpf/src/profiler
    )
    add_library(profiler SHARED IMPORTED)
    set_target_properties(profiler PROPERTIES
        IMPORTED_LOCATION ${CMAKE_BINARY_DIR}/_thirdparty/coolbpf/src/profiler/release/libprofiler.so
        IMPORTED_NO_SONAME TRUE
    )
    add_dependencies(profiler build_libprofiler)
    add_custom_target(
        install_coolbpf ALL
        COMMAND ${CMAKE_COMMAND} --install ${CMAKE_BINARY_DIR}/_thirdparty/coolbpf
        DEPENDS ${_coolbpf_install_deps} build_libprofiler
    )

    include_directories("${DEPS_INCLUDE_ROOT}/coolbpf")

    file(GLOB LIB_SOURCE_FILES *.cpp *.h)
    if (BUILD_LOGTAIL_UT)
        add_definitions(-DAPSARA_UNIT_TEST_MAIN)
    endif ()

    add_library(${EBPF_DIRVER_TARGET} SHARED ${LIB_SOURCE_FILES})
    target_link_options(${EBPF_DIRVER_TARGET} PRIVATE "-Wl,-rpath=.:${DEFAULT_DEPS_ROOT}/lib")
    target_link_libraries(${EBPF_DIRVER_TARGET} coolbpf profiler)
    add_dependencies(${EBPF_DIRVER_TARGET} install_coolbpf)
    target_compile_options(${EBPF_DIRVER_TARGET} PRIVATE -fPIC)
endif ()
