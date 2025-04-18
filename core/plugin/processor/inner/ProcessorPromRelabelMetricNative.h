/*
 * Copyright 2024 iLogtail Authors
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

#include "collection_pipeline/plugin/interface/Processor.h"
#include "models/PipelineEventGroup.h"
#include "models/PipelineEventPtr.h"
#include "prometheus/schedulers/ScrapeConfig.h"

namespace logtail {

namespace prom {
struct AutoMetric {
    double mScrapeDurationSeconds;
    uint64_t mScrapeResponseSizeBytes;
    uint64_t mScrapeSamplesLimit;
    // uint64_t mPostRelabel;
    uint64_t mScrapeSamplesScraped;
    uint64_t mScrapeTimeoutSeconds;
    std::string mScrapeState;
    bool mUp;
};
} // namespace prom

class ProcessorPromRelabelMetricNative : public Processor {
public:
    static const std::string sName;

    const std::string& Name() const override { return sName; }
    bool Init(const Json::Value& config) override;
    void Process(PipelineEventGroup& metricGroup) override;

protected:
    bool IsSupportedEvent(const PipelineEventPtr& e) const override;

private:
    bool ProcessEvent(PipelineEventPtr& e, const GroupTags& targetTags);

    void AddAutoMetrics(PipelineEventGroup& eGroup, const prom::AutoMetric& autoMetric) const;
    void UpdateAutoMetrics(const PipelineEventGroup& eGroup, prom::AutoMetric& autoMetric) const;
    void AddMetric(PipelineEventGroup& metricGroup,
                   const std::string& name,
                   double value,
                   time_t timestamp,
                   uint32_t nanoSec,
                   const GroupTags& targetTags) const;

    std::unique_ptr<ScrapeConfig> mScrapeConfigPtr;
    std::string mLoongCollectorScraper;

#ifdef APSARA_UNIT_TEST_MAIN
    friend class ProcessorPromRelabelMetricNativeUnittest;
    friend class InputPrometheusUnittest;
#endif
};

} // namespace logtail
