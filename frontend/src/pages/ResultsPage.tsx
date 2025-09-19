import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { MetricDial } from "@/components/ui/metric-dial";
import { Chatbot } from "@/components/ui/chatbot";
import { ArrowLeft, Download, Share, AlertTriangle, TrendingUp, Heart, Activity, Zap, Loader2 } from "lucide-react";
import { useNavigate, useParams } from "react-router-dom";
import { reportsApi, Report, HealthMetric } from "@/lib/api";

interface ReportSummary extends Report {
  summary: string;
}

const ResultsPage = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [summary, setSummary] = useState<ReportSummary | null>(null);
  const [metrics, setMetrics] = useState<HealthMetric[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchReportData = async () => {
      console.log("Fetching data for report ID:", id);
      if (!id) {
        setError("Report ID not found.");
        setLoading(false);
        return;
      }
      try {
        setLoading(true);
        const reportId = parseInt(id, 10);
        console.log("Parsed report ID:", reportId);

        const summaryData = await reportsApi.getSummary(reportId);
        console.log("Summary data:", summaryData);

        const metricsData = await reportsApi.getHealthMetrics(reportId);
        console.log("Metrics data:", metricsData);

        setSummary({ ...summaryData.report, summary: summaryData.summary });
        setMetrics(metricsData.metrics);
        console.log("State updated:", { summary, metrics });

      } catch (err: any) {
        console.error("Failed to fetch report data:", err);
        setError(err.message || "Failed to fetch report data.");
      } finally {
        setLoading(false);
      }
    };

    fetchReportData();
  }, [id]);

  const getStatusColor = (status: string) => {
    switch (status) {
      case "normal":
        return "health-normal";
      case "high":
      case "low":
        return status === "high" ? "health-danger" : "health-warning";
      default:
        return "muted";
    }
  };

  const getStatusBg = (status: string) => {
    switch (status) {
      case "normal":
        return "health-normal-bg";
      case "high":
      case "low":
        return status === "high" ? "health-danger-bg" : "health-warning-bg";
      default:
        return "muted";
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "normal":
        return { text: "Normal", variant: "default" as const };
      case "high":
        return { text: "High", variant: "destructive" as const };
      case "low":
        return { text: "Low", variant: "secondary" as const };
      default:
        return { text: "Unknown", variant: "outline" as const };
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="flex items-center space-x-2">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
          <span className="text-lg text-muted-foreground">Analyzing your report...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center space-y-4">
          <AlertTriangle className="h-12 w-12 text-destructive mx-auto" />
          <h2 className="text-2xl font-semibold text-foreground">Error</h2>
          <p className="text-muted-foreground">{error}</p>
          <Button onClick={() => navigate("/analyze")}>Go Back</Button>
        </div>
      </div>
    );
  }

  if (!summary) {
    return null;
  }

  const parsedSummary = JSON.parse(summary.summary);

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="p-6 border-b border-border/50">
        <div className="max-w-4xl mx-auto flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <Button 
              variant="ghost" 
              size="sm" 
              onClick={() => navigate("/analyze")}
              className="transition-smooth"
            >
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back
            </Button>
            <h1 className="text-2xl font-semibold text-foreground">Your Report Summary</h1>
          </div>
          
          <div className="flex items-center space-x-2">
            <Button variant="outline" size="sm" className="transition-smooth">
              <Share className="h-4 w-4 mr-2" />
              Share
            </Button>
            <Button variant="outline" size="sm" className="transition-smooth">
              <Download className="h-4 w-4 mr-2" />
              Download
            </Button>
          </div>
        </div>
      </header>

      <main className="p-6">
        <div className="max-w-4xl mx-auto space-y-6">
          {/* Disclaimer */}
          <Card className="p-6 bg-amber-50/50 border-amber-200">
            <div className="flex items-start space-x-3">
              <AlertTriangle className="h-5 w-5 text-amber-600 mt-0.5 flex-shrink-0" />
              <div className="space-y-1">
                <h3 className="font-medium text-amber-900">Important Notice</h3>
                <p className="text-sm text-amber-800">
                  This is an AI-generated summary and not medical advice. 
                  Please consult your doctor to discuss your results and any concerns.
                </p>
              </div>
            </div>
          </Card>

          {/* Metric Dials Overview */}
          <Card className="p-6">
            <h2 className="text-xl font-semibold text-foreground mb-6 flex items-center">
              <Activity className="h-5 w-5 mr-2 text-primary" />
              Your Health Metrics
            </h2>
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-6">
              {metrics.map((metric, index) => (
                <MetricDial
                  key={index}
                  title={metric.name}
                  value={`${metric.value} ${metric.unit}`}
                  normalRange={`${metric.range_min}-${metric.range_max} ${metric.unit}`}
                  percentage={metric.score}
                  status={metric.status as "normal" | "warning" | "danger"}
                />
              ))}
            </div>
          </Card>

          {/* AI Analysis Overview */}
          <Card className="p-6 bg-gradient-card">
            <h2 className="text-xl font-semibold text-foreground mb-4 flex items-center">
              <Zap className="h-5 w-5 mr-2 text-primary" />
              AI Analysis Summary
            </h2>
            
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Overall Health */}
              <div className="space-y-4">
                <div className="flex items-center space-x-2">
                  <Heart className="h-4 w-4 text-health-normal" />
                  <span className="font-medium text-foreground">Overall Health Status</span>
                </div>
                <div className="p-4 bg-health-normal-bg rounded-xl">
                  <span className="text-lg font-semibold text-health-normal">{parsedSummary.overall_health_status}</span>
                </div>
              </div>

              {/* Key Findings */}
              <div className="space-y-4">
                <div className="flex items-center space-x-2">
                  <TrendingUp className="h-4 w-4 text-primary" />
                  <span className="font-medium text-foreground">Key Findings</span>
                </div>
                <div className="space-y-2">
                  {parsedSummary.key_findings.map((finding: string, index: number) => (
                    <div key={index} className="p-3 bg-primary-soft rounded-lg text-sm text-foreground">
                      • {finding}
                    </div>
                  ))}
                </div>
              </div>
            </div>

            {/* Recommendations */}
            <div className="mt-6 p-4 bg-secondary/20 rounded-xl">
              <h3 className="font-semibold text-foreground mb-3">Recommended Actions</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                {parsedSummary.recommendations.map((rec: string, index: number) => (
                  <div key={index} className="flex items-start space-x-2 text-sm">
                    <span className="text-secondary-accent mt-1">✓</span>
                    <span className="text-foreground">{rec}</span>
                  </div>
                ))}
              </div>
            </div>

            {/* Next Steps */}
            <div className="mt-4 p-4 bg-health-warning-bg rounded-xl">
              <h3 className="font-semibold text-health-warning mb-2">Next Steps</h3>
              <p className="text-sm text-foreground">{parsedSummary.next_steps}</p>
            </div>
          </Card>

          {/* Results */}
          <div className="space-y-4">
            {metrics.map((metric, index) => (
              <Card key={index} className="p-6 transition-smooth hover:shadow-card">
                <div className="flex items-start justify-between mb-4">
                  <div>
                    <h3 className="text-lg font-semibold text-foreground">{metric.name}</h3>
                    <div className="flex items-center space-x-3 mt-1">
                      <span className="text-2xl font-bold text-foreground">{`${metric.value} ${metric.unit}`}</span>
                      <Badge 
                        variant={getStatusBadge(metric.status).variant}
                        className="transition-smooth"
                      >
                        {getStatusBadge(metric.status).text}
                      </Badge>.
                    </div>
                  </div>
                  <div className={`w-1 h-16 rounded-full bg-${getStatusColor(metric.status)}`} />
                </div>
                
                <div className="space-y-3">
                  <div className="text-sm text-muted-foreground">
                    <span className="font-medium">Normal Range:</span> {`${metric.range_min}-${metric.range_max} ${metric.unit}`}
                  </div>
                  
                  <div className={`p-4 rounded-xl bg-${getStatusBg(metric.status)}`}>
                    <h4 className="font-medium text-foreground mb-2">What this means:</h4>
                    <p className="text-sm text-foreground leading-relaxed">
                      {metric.description}
                    </p>
                  </div>
                </div>
              </Card>
            ))}
          </div>

          {/* Actions */}
          <Card className="p-6 text-center">
            <h3 className="text-lg font-semibold text-foreground mb-2">Need Help Understanding?</h3>
            <p className="text-muted-foreground mb-4">
              Our AI assistant can answer specific questions about your report.
            </p>
            <Button className="transition-smooth">
              Ask Questions About My Report
            </Button>
          </Card>
        </div>
      </main>

      {/* Chatbot */}
      <Chatbot />
    </div>
  );
};

export default ResultsPage;
