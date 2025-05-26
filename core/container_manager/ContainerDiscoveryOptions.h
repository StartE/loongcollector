/*
 * Copyright 2023 iLogtail Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#pragma once

#include <string>
#include <unordered_map>
#include <utility>

#include "json/json.h"

#include "collection_pipeline/CollectionPipelineContext.h"
#include "collection_pipeline/plugin/instance/PluginInstance.h"
#include <boost/regex.hpp>

namespace logtail {


struct FieldFilter {
    std::unordered_map<std::string, std::string> mFieldsMap;
    std::unordered_map<std::string, std::shared_ptr<boost::regex>> mFieldsRegMap;
};

struct MatchCriteriaFilter {
    // 包含和排除的标签
    FieldFilter mIncludeFields;
    FieldFilter mExcludeFields;
};

struct K8sFilter {
    std::shared_ptr<boost::regex> mNamespaceReg;
    std::shared_ptr<boost::regex> mPodReg;
    std::shared_ptr<boost::regex> mContainerReg;

    MatchCriteriaFilter mK8sLabelFilter;
};

struct ContainerFilters {
    K8sFilter mK8SFilter;
    MatchCriteriaFilter mEnvFilter;    
    MatchCriteriaFilter mContainerLabelFilter;
};

struct ContainerFilterConfig {
    std::string mK8sNamespaceRegex;
    std::string mK8sPodRegex;
    std::string mK8sContainerRegex;


    std::unordered_map<std::string, std::string> mIncludeK8sLabel;
    std::unordered_map<std::string, std::string> mExcludeK8sLabel;

    std::unordered_map<std::string, std::string> mIncludeEnv;
    std::unordered_map<std::string, std::string> mExcludeEnv;

    std::unordered_map<std::string, std::string> mIncludeContainerLabel;
    std::unordered_map<std::string, std::string> mExcludeContainerLabel;

    bool Init(const Json::Value& config, const CollectionPipelineContext& ctx, const std::string& pluginType);
    
    bool GetContainerFilters(ContainerFilters& mContainerFilters);
};



struct ContainerDiscoveryOptions {
    ContainerFilterConfig mContainerFilterConfig;
    ContainerFilters mContainerFilters;
    std::unordered_map<std::string, std::string> mExternalK8sLabelTag;
    std::unordered_map<std::string, std::string> mExternalEnvTag;
    // 启用容器元信息预览
    bool mCollectingContainersMeta = false;

    bool Init(const Json::Value& config, const CollectionPipelineContext& ctx, const std::string& pluginType);
    /*
    void GenerateContainerMetaFetchingGoPipeline(Json::Value& res,
                                                 const FileDiscoveryOptions* fileDiscovery = nullptr,
                                                 const PluginInstance::PluginMeta& pluginMeta = {"0"}) const;
    */
};

using ContainerDiscoveryConfig = std::pair<const ContainerDiscoveryOptions*, const CollectionPipelineContext*>;

} // namespace logtail
