import { useState, useCallback } from "react";
import { useDropzone } from "react-dropzone";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Upload, FileText, ArrowLeft, Loader2 } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { useToast } from "@/hooks/use-toast";

const AnalyzePage = () => {
  const navigate = useNavigate();
  const { toast } = useToast();
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [progress, setProgress] = useState(0);
  const [uploadedFile, setUploadedFile] = useState<File | null>(null);

  const onDrop = useCallback((acceptedFiles: File[]) => {
    if (acceptedFiles.length > 0) {
      setUploadedFile(acceptedFiles[0]);
      toast({
        title: "File uploaded successfully",
        description: `${acceptedFiles[0].name} is ready for analysis.`,
        duration: 3000,
      });
    }
  }, [toast]);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      'application/pdf': ['.pdf'],
      'text/plain': ['.txt'],
      'image/*': ['.png', '.jpg', '.jpeg']
    },
    maxFiles: 1,
    multiple: false
  });

  const handleAnalyze = async () => {
    if (!uploadedFile) return;

    const token = localStorage.getItem("token");
    if (!token) {
      toast({
        title: "Authentication error",
        description: "You must be logged in to analyze a report.",
        variant: "destructive",
      });
      navigate("/login");
      return;
    }

    setIsAnalyzing(true);
    setProgress(0);

    const formData = new FormData();
    formData.append("file", uploadedFile);

    try {
      const response = await fetch("/api/reports/upload", { // Replace with your actual upload endpoint
        method: "POST",
        headers: {
          "Authorization": `Bearer ${token}`,
        },
        body: formData,
      });

      if (response.ok) {
        // Simulate analysis progress
        const progressInterval = setInterval(() => {
          setProgress(prev => {
            if (prev >= 90) {
              clearInterval(progressInterval);
              return 90;
            }
            return prev + 10;
          });
        }, 300);

        // Simulate API call to get analysis results
        setTimeout(() => {
          clearInterval(progressInterval);
          setProgress(100);

          setTimeout(() => {
            navigate("/results");
          }, 500);
        }, 3000);
      } else {
        const errorData = await response.json();
        toast({
          title: "Analysis failed",
          description: errorData.error || "An unknown error occurred.",
          variant: "destructive",
        });
        setIsAnalyzing(false);
      }
    } catch (error) {
      toast({
        title: "Analysis failed",
        description: "An unknown error occurred.",
        variant: "destructive",
      });
      setIsAnalyzing(false);
    }
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="p-6 border-b border-border/50">
        <div className="max-w-4xl mx-auto flex items-center space-x-4">
          <Button 
            variant="ghost" 
            size="sm" 
            onClick={() => navigate("/")}
            className="transition-smooth"
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back
          </Button>
          <h1 className="text-2xl font-semibold text-foreground">Analyze Your Report</h1>
        </div>
      </header>

      <main className="p-6">
        <div className="max-w-4xl mx-auto">
          {!isAnalyzing ? (
            <div className="space-y-8">
              {/* Enhanced Upload area */}
              <Card className="p-8 relative overflow-hidden">
                {/* Background decoration */}
                <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-secondary/5" />
                
                <div
                  {...getRootProps()}
                  className={`relative border-2 border-dashed rounded-2xl p-12 text-center cursor-pointer transition-bounce ${
                    isDragActive 
                      ? 'border-primary bg-primary/10 scale-105' 
                      : 'border-border hover:border-primary/50 hover:bg-primary/5 hover:scale-102'
                  }`}
                >
                  <input {...getInputProps()} />
                  
                  {/* Floating particles effect */}
                  <div className="absolute inset-0 pointer-events-none">
                    <div className="absolute top-4 left-8 w-2 h-2 bg-primary/20 rounded-full animate-pulse" />
                    <div className="absolute top-8 right-12 w-1 h-1 bg-secondary-accent/30 rounded-full animate-pulse" style={{animationDelay: '0.5s'}} />
                    <div className="absolute bottom-6 left-16 w-1.5 h-1.5 bg-health-normal/25 rounded-full animate-pulse" style={{animationDelay: '1s'}} />
                    <div className="absolute bottom-12 right-6 w-1 h-1 bg-primary/15 rounded-full animate-pulse" style={{animationDelay: '1.5s'}} />
                  </div>
                  
                  <div className="space-y-6">
                    <div className={`w-20 h-20 bg-primary/10 rounded-3xl flex items-center justify-center mx-auto transition-bounce ${
                      isDragActive ? 'scale-110 bg-primary/20' : 'hover:scale-105'
                    }`}>
                      <Upload className={`h-10 w-10 text-primary transition-smooth ${
                        isDragActive ? 'animate-bounce' : ''
                      }`} />
                    </div>
                    
                    {uploadedFile ? (
                      <div className="space-y-3 animate-fade-in">
                        <div className="flex items-center justify-center space-x-3 text-health-normal">
                          <div className="p-2 bg-health-normal/10 rounded-full">
                            <FileText className="h-6 w-6" />
                          </div>
                          <span className="font-semibold text-lg">{uploadedFile.name}</span>
                        </div>
                        <div className="flex items-center justify-center space-x-2">
                          <div className="w-2 h-2 bg-health-normal rounded-full animate-pulse" />
                          <p className="text-sm text-muted-foreground font-medium">
                            File uploaded successfully. Ready to analyze!
                          </p>
                          <div className="w-2 h-2 bg-health-normal rounded-full animate-pulse" />
                        </div>
                      </div>
                    ) : (
                      <div className="space-y-4">
                        <div className="space-y-2">
                          <h3 className="text-2xl font-bold text-foreground">
                            {isDragActive ? (
                              <span className="text-primary animate-pulse">
                                🎯 Drop your file here
                              </span>
                            ) : (
                              "📋 Upload your medical report"
                            )}
                          </h3>
                          <p className="text-muted-foreground text-lg">
                            Drag and drop your report, or click to browse
                          </p>
                        </div>
                        
                        {/* Feature highlights */}
                        <div className="grid grid-cols-3 gap-4 mt-6">
                          <div className="p-3 bg-primary/5 rounded-xl">
                            <div className="text-2xl mb-1">🔒</div>
                            <div className="text-xs font-medium text-primary">Secure</div>
                          </div>
                          <div className="p-3 bg-secondary/20 rounded-xl">
                            <div className="text-2xl mb-1">⚡</div>
                            <div className="text-xs font-medium text-secondary-accent">Fast</div>
                          </div>
                          <div className="p-3 bg-health-normal/10 rounded-xl">
                            <div className="text-2xl mb-1">🎯</div>
                            <div className="text-xs font-medium text-health-normal">Accurate</div>
                          </div>
                        </div>
                        
                        <div className="text-sm text-muted-foreground bg-muted/50 rounded-lg p-3">
                          📄 Supports: PDF, TXT, or image files (up to 20MB)
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              </Card>

              {/* Enhanced Action button */}
              {uploadedFile && (
                <div className="text-center animate-fade-in">
                  <div className="space-y-4">
                    <div className="flex items-center justify-center space-x-2 text-muted-foreground">
                      <div className="h-px bg-border flex-1" />
                      <span className="text-sm font-medium">Ready to analyze</span>
                      <div className="h-px bg-border flex-1" />
                    </div>
                    
                    <Button 
                      size="lg"
                      onClick={handleAnalyze}
                      className="px-12 py-6 text-lg font-semibold shadow-card hover:shadow-floating transition-bounce hover:scale-105 bg-gradient-to-r from-primary to-secondary-accent hover:from-primary/90 hover:to-secondary-accent/90"
                    >
                      🔬 Analyze My Report
                      <ArrowLeft className="ml-3 h-6 w-6 rotate-180" />
                    </Button>
                    
                    <p className="text-sm text-muted-foreground">
                      Our AI will analyze your report in under 30 seconds
                    </p>
                  </div>
                </div>
              )}
            </div>
          ) : (
            /* Analysis in progress */
            <Card className="p-12 text-center space-y-6">
              <div className="w-20 h-20 bg-primary/10 rounded-full flex items-center justify-center mx-auto">
                <Loader2 className="h-10 w-10 text-primary animate-spin" />
              </div>
              
              <div className="space-y-4">
                <h2 className="text-2xl font-semibold text-foreground">
                  Analyzing your results...
                </h2>
                <p className="text-muted-foreground max-w-md mx-auto">
                  Our AI is carefully reviewing your medical report and translating 
                  complex terms into simple, understandable language.
                </p>
              </div>
              
              <div className="max-w-md mx-auto space-y-2">
                <Progress value={progress} className="h-2" />
                <div className="text-sm text-muted-foreground">
                  {progress < 30 && "Reading your report..."}
                  {progress >= 30 && progress < 60 && "Identifying key biomarkers..."}
                  {progress >= 60 && progress < 90 && "Generating explanations..."}
                  {progress >= 90 && "Almost done!"}
                </div>
              </div>
            </Card>
          )}
        </div>
      </main>
    </div>
  );
};

export default AnalyzePage;
