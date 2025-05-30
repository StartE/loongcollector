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

#include "TagConstants.h"
#include "collection_pipeline/plugin/interface/Processor.h"
#include "common/StringView.h"

namespace logtail {

class ProcessorTagNative : public Processor {
public:
    static const std::string sName;

    const std::string& Name() const override { return sName; }
    bool Init(const Json::Value& config) override;
    void Process(PipelineEventGroup& logGroup) override;

protected:
    bool IsSupportedEvent(const PipelineEventPtr& e) const override;

private:
    void AddTag(PipelineEventGroup& logGroup, TagKey tagKey, const std::string& value) const;
    void AddTag(PipelineEventGroup& logGroup, TagKey tagKey, StringView value) const;
    std::unordered_map<TagKey, std::string> mPipelineMetaTagKey;
    // After unmarshalling from json, we cannot determine the map is empty or no such config
    bool mAppendingAllEnvMetaTag = false;
    std::unordered_map<std::string, std::string> mAgentEnvMetaTagKey;

#ifdef APSARA_UNIT_TEST_MAIN
    friend class ProcessorTagNativeUnittest;
#endif
};

} // namespace logtail
