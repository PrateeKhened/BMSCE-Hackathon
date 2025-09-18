import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { MetricDial } from "@/components/ui/metric-dial";
import { Chatbot } from "@/components/ui/chatbot";
import { ArrowLeft, Download, Share, AlertTriangle, TrendingUp, Heart, Activity, Zap } from "lucide-react";
import { useNavigate } from "react-router-dom";

// Mock data for demonstration
const mockResults = [
  {
    biomarker: "Hemoglobin",
    value: "11.5 g/dL",
    normalRange: "13.5 - 17.5 g/dL",
    status: "low",
    percentage: 67,
    explanation: "Your hemoglobin is slightly below the normal range. This can sometimes be a sign of anemia, which means your body may not be getting enough oxygen. This is common and often treatable."
  },
  {
    biomarker: "LDL Cholesterol",
    value: "135 mg/dL",
    normalRange: "< 100 mg/dL",
    status: "high",
    percentage: 135,
    explanation: "Your LDL (\"bad\") cholesterol is elevated. This increases your risk of heart disease. Simple dietary changes and exercise can often help bring this number down."
  },
  {
    biomarker: "Blood Glucose",
    value: "92 mg/dL",
    normalRange: "70 - 100 mg/dL",
    status: "normal",
    percentage: 92,
    explanation: "Your blood sugar level is in the healthy range. This indicates good blood sugar control and a lower risk of diabetes."
  },
  {
    biomarker: "White Blood Cells",
    value: "7,200 cells/μL",
    normalRange: "4,000 - 11,000 cells/μL",
    status: "normal",
    percentage: 85,
    explanation: "Your white blood cell count is normal, indicating your immune system is functioning well and there are no signs of infection or inflammation."
  }
];

const mockAiAnalysis = {
  overallHealth: "Generally Good",
  keyFindings: [
    "Slight anemia indicated by low hemoglobin levels",
    "Elevated LDL cholesterol requiring dietary attention",
    "Normal blood sugar and immune function"
  ],
  recommendations: [
    "Consider iron-rich foods to boost hemoglobin",
    "Reduce saturated fats to lower LDL cholesterol",
    "Maintain current healthy lifestyle for blood sugar",
    "Schedule follow-up in 3 months"
  ],
  riskFactors: ["Cardiovascular risk due to high LDL"],
  nextSteps: "Consult with your doctor about dietary changes and possible iron supplementation."
};

const ResultsPage = () => {
  const navigate = useNavigate();

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
              {mockResults.map((result, index) => (
                <MetricDial
                  key={index}
                  title={result.biomarker}
                  value={result.value}
                  normalRange={result.normalRange}
                  percentage={result.percentage}
                  status={result.status as "normal" | "warning" | "danger"}
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
                  <span className="text-lg font-semibold text-health-normal">{mockAiAnalysis.overallHealth}</span>
                </div>
              </div>

              {/* Key Findings */}
              <div className="space-y-4">
                <div className="flex items-center space-x-2">
                  <TrendingUp className="h-4 w-4 text-primary" />
                  <span className="font-medium text-foreground">Key Findings</span>
                </div>
                <div className="space-y-2">
                  {mockAiAnalysis.keyFindings.map((finding, index) => (
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
                {mockAiAnalysis.recommendations.map((rec, index) => (
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
              <p className="text-sm text-foreground">{mockAiAnalysis.nextSteps}</p>
            </div>
          </Card>

          {/* Results */}
          <div className="space-y-4">
            {mockResults.map((result, index) => (
              <Card key={index} className="p-6 transition-smooth hover:shadow-card">
                <div className="flex items-start justify-between mb-4">
                  <div>
                    <h3 className="text-lg font-semibold text-foreground">{result.biomarker}</h3>
                    <div className="flex items-center space-x-3 mt-1">
                      <span className="text-2xl font-bold text-foreground">{result.value}</span>
                      <Badge 
                        variant={getStatusBadge(result.status).variant}
                        className="transition-smooth"
                      >
                        {getStatusBadge(result.status).text}
                      </Badge>
                    </div>
                  </div>
                  <div className={`w-1 h-16 rounded-full bg-${getStatusColor(result.status)}`} />
                </div>
                
                <div className="space-y-3">
                  <div className="text-sm text-muted-foreground">
                    <span className="font-medium">Normal Range:</span> {result.normalRange}
                  </div>
                  
                  <div className={`p-4 rounded-xl bg-${getStatusBg(result.status)}`}>
                    <h4 className="font-medium text-foreground mb-2">What this means:</h4>
                    <p className="text-sm text-foreground leading-relaxed">
                      {result.explanation}
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