/*
 * Copyright 2023 iLogtail Authors
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

#include <memory>
#include <set>
#include <string>
#include <unordered_map>

#include "collection_pipeline/plugin/creator/PluginCreator.h"
#include "collection_pipeline/plugin/instance/FlusherInstance.h"
#include "collection_pipeline/plugin/instance/InputInstance.h"
#include "collection_pipeline/plugin/instance/PluginInstance.h"
#include "collection_pipeline/plugin/instance/ProcessorInstance.h"
#include "common/DynamicLibHelper.h"
#include "runner/sink/SinkType.h"

struct processor_interface_t;

namespace logtail {

class PluginRegistry {
public:
    PluginRegistry(const PluginRegistry&) = delete;
    PluginRegistry& operator=(const PluginRegistry&) = delete;

    static PluginRegistry* GetInstance() {
        static PluginRegistry instance;
        return &instance;
    }

    void LoadPlugins();
    void UnloadPlugins();
    std::unique_ptr<InputInstance> CreateInput(const std::string& name, const PluginInstance::PluginMeta& pluginMeta);
    std::unique_ptr<ProcessorInstance> CreateProcessor(const std::string& name,
                                                       const PluginInstance::PluginMeta& pluginMeta);
    std::unique_ptr<FlusherInstance> CreateFlusher(const std::string& name,
                                                   const PluginInstance::PluginMeta& pluginMeta);
    bool IsValidGoPlugin(const std::string& name) const;
    bool IsValidNativeInputPlugin(const std::string& name) const;
    bool IsValidNativeProcessorPlugin(const std::string& name) const;
    bool IsValidNativeFlusherPlugin(const std::string& name) const;
    bool IsGlobalSingletonInputPlugin(const std::string& name) const;

private:
    enum PluginCat { INPUT_PLUGIN, PROCESSOR_PLUGIN, FLUSHER_PLUGIN };

    struct PluginKey {
        PluginCat cat;
        std::string name;
        PluginKey(PluginCat cat, const std::string& name) : cat(cat), name(name) {}
        bool operator==(const PluginKey& other) const { return cat == other.cat && name == other.name; }
    };

    struct PluginKeyHash {
        size_t operator()(const PluginKey& obj) const {
            return std::hash<int>()(obj.cat) ^ std::hash<std::string>()(obj.name);
        }
    };

    using PluginCreatorWithInfo = std::pair<std::shared_ptr<PluginCreator>, bool>;

    PluginRegistry() {}
    ~PluginRegistry() = default;

    void LoadStaticPlugins();
    void LoadDynamicPlugins(const std::set<std::string>& plugins);
    void RegisterInputCreator(PluginCreator* creator, bool isSingleton = false);
    void RegisterProcessorCreator(PluginCreator* creator);
    void RegisterFlusherCreator(PluginCreator* creator, bool isSingleton = false);
    PluginCreator* LoadProcessorPlugin(DynamicLibLoader& loader, const std::string pluginType);
    void RegisterCreator(PluginCat cat, PluginCreator* creator, bool isSingleton);
    std::unique_ptr<PluginInstance>
    Create(PluginCat cat, const std::string& name, const PluginInstance::PluginMeta& pluginMeta);
    bool IsGlobalSingleton(PluginCat cat, const std::string& name) const;

    std::unordered_map<PluginKey, PluginCreatorWithInfo, PluginKeyHash> mPluginDict;

#ifdef APSARA_UNIT_TEST_MAIN
    friend class PluginRegistryUnittest;
    friend class PipelineUpdateUnittest;
    friend void LoadPluginMock();
#endif
};

} // namespace logtail
