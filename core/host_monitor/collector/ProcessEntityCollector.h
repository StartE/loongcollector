/*
 * Copyright 2024 iLogtail Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#pragma once

#include <cstdint>
#include <cstdlib>
#include <sys/types.h>
#include <unistd.h>

#include <boost/filesystem.hpp>
#include <filesystem>
#include <memory>
#include <string>
#include <unordered_map>

#include "constants/EntityConstants.h"
#include "host_monitor/collector/BaseCollector.h"
#include "models/StringView.h"

using namespace std::chrono;

namespace logtail {

struct ProcessCpuInfo {
    uint64_t user = 0;
    uint64_t sys = 0;
    uint64_t total = 0;
    double percent = 0.0;
};

struct ProcessStat {
    pid_t pid = 0;
    uint64_t vSize = 0;
    uint64_t rss = 0;
    uint64_t minorFaults = 0;
    uint64_t majorFaults = 0;
    pid_t parentPid = 0;
    int tty = 0;
    int priority = 0;
    int nice = 0;
    int numThreads = 0;
    system_clock::time_point startTime;
    steady_clock::time_point lastTime;

    milliseconds utime{0};
    milliseconds stime{0};
    milliseconds cutime{0};
    milliseconds cstime{0};
    ProcessCpuInfo cpuInfo;

    std::string name;
    char state = '\0';
    int processor = 0;

    int64_t startMillis() const { return duration_cast<milliseconds>(startTime.time_since_epoch()).count(); }
};

typedef std::shared_ptr<ProcessStat> ProcessStatPtr;

// See https://man7.org/linux/man-pages/man5/proc.5.html
enum class EnumProcessStat : int {
    pid, // 0
    comm, // 1
    state, // 2
    ppid, // 3
    pgrp, // 4
    session, // 5
    tty_nr, // 6
    tpgid, // 7
    flags, // 8
    minflt, // 9
    cminflt, // 10
    majflt, // 11
    cmajflt, // 12
    utime, // 13
    stime, // 14
    cutime, // 15
    cstime, // 16
    priority, // 17
    nice, // 18
    num_threads, // 19
    itrealvalue, // 20
    starttime, // 21
    vsize, // 22
    rss, // 23
    rsslim, // 24
    startcode, // 25
    endcode, // 26
    startstack, // 27
    kstkesp, // 28
    kstkeip, // 29
    signal, // 30
    blocked, // 31
    sigignore, // 32
    sigcatch, // 33
    wchan, // 34
    nswap, // 35
    cnswap, // 36
    exit_signal, // 37
    processor, // 38 <--- 至少需要有该字段
    rt_priority, // 39
    policy, // 40
    delayacct_blkio_ticks, // 41
    guest_time, // 42
    cguest_time, // 43
    start_data, // 44
    end_data, // 45
    start_brk, // 46
    arg_start, // 47
    arg_end, // 48
    env_start, // 49
    env_end, // 50
    exit_code, // 51

    _count, // 只是用于计数，非实际字段
};
static_assert((int)EnumProcessStat::comm == 1, "EnumProcessStat invalid");
static_assert((int)EnumProcessStat::processor == 38, "EnumProcessStat invalid");

constexpr int operator-(EnumProcessStat a, EnumProcessStat b) {
    return (int)a - (int)b;
}

class ProcessEntityCollector : public BaseCollector {
public:
    ProcessEntityCollector();
    ~ProcessEntityCollector() override = default;

    bool Collect(const HostMonitorTimerEvent::CollectConfig& collectConfig, PipelineEventGroup* group) override;

    static const std::string sName;
    const std::string& Name() const override { return sName; }

private:
    void GetSortedProcess(std::vector<ProcessStatPtr>& processStats, size_t topN);
    ProcessStatPtr GetProcessStat(pid_t pid, bool& isFirstCollect);
    ProcessStatPtr ReadNewProcessStat(pid_t pid);
    ProcessStatPtr ParseProcessStat(pid_t pid, std::string& line);
    bool WalkAllProcess(const std::filesystem::path& root, const std::function<void(const std::string&)>& callback);

    const std::string GetProcessEntityID(StringView pid, StringView createTime, StringView hostEntityID);
    void FetchDomainInfo(std::string& domain,
                         std::string& entityType,
                         std::string& hostEntityType,
                         StringView& hostEntityID);
    int64_t GetHostSystemBootTime();

    steady_clock::time_point mProcessSortTime;
    std::unordered_map<pid_t, ProcessStatPtr> mPrevProcessStat;

    const int mProcessSilentCount;

#ifdef APSARA_UNIT_TEST_MAIN
    friend class ProcessEntityCollectorUnittest;
#endif
};

} // namespace logtail
