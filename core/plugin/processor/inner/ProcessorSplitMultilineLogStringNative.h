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

#include <cstdint>

#include <vector>

#include "collection_pipeline/plugin/interface/Processor.h"
#include "constants/Constants.h"
#include "file_server/MultilineOptions.h"
#include "plugin/processor/CommonParserOptions.h"

namespace logtail {

class ProcessorSplitMultilineLogStringNative : public Processor {
public:
    static const std::string sName;

    std::string mSourceKey = DEFAULT_CONTENT_KEY;
    MultilineOptions mMultiline;
    bool mEnableRawContent = false;

    const std::string& Name() const override { return sName; }
    bool Init(const Json::Value& config) override;
    void Process(PipelineEventGroup& logGroup) override;

protected:
    bool IsSupportedEvent(const PipelineEventPtr& e) const override;

private:
    void ProcessEvent(PipelineEventGroup& logGroup,
                      StringView logPath,
                      PipelineEventPtr&& e,
                      EventsContainer& newEvents,
                      int* inputLines,
                      int* unmatchLines);
    void CreateNewEvent(const StringView& content,
                        bool isLastLog,
                        StringBuffer& sourceKey,
                        const LogEvent& sourceEvent,
                        PipelineEventGroup& logGroup,
                        EventsContainer& newEvents);
    void HandleUnmatchLogs(const StringView& sourceVal,
                           bool isLastLog,
                           StringBuffer& sourceKey,
                           const LogEvent& sourceEvent,
                           PipelineEventGroup& logGroup,
                           EventsContainer& newEvents,
                           StringView logPath,
                           int* unmatchLines);
    StringView GetNextLine(StringView log, size_t begin);

    bool HasStartPattern() const { return !mStartPatternReg.empty(); }
    bool HasContinuePattern() const { return !mContinuePatternReg.empty(); }
    bool HasEndPattern() const { return !mEndPatternReg.empty(); }
    const boost::regex& GetStartPatternReg() const;
    const boost::regex& GetContinuePatternReg() const;
    const boost::regex& GetEndPatternReg() const;

    // boost::regex object shared by multi-thread leads to performance degradation. Therefore, each thread should be
    // allocated a different copy.
    std::vector<boost::regex> mStartPatternReg;
    std::vector<boost::regex> mContinuePatternReg;
    std::vector<boost::regex> mEndPatternReg;

    CounterPtr mMatchedEventsTotal;
    CounterPtr mMatchedLinesTotal;
    CounterPtr mUnmatchedLinesTotal;

#ifdef APSARA_UNIT_TEST_MAIN
    friend class ProcessorSplitMultilineLogStringNativeUnittest;
    friend class ProcessorSplitMultilineLogDisacardUnmatchUnittest;
    friend class ProcessorSplitMultilineLogKeepUnmatchUnittest;
#endif
};

} // namespace logtail
