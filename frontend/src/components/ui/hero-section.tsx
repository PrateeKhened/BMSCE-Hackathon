import { Button } from "@/components/ui/button";
import { ArrowRight, Shield, Zap, Heart } from "lucide-react";
import heroImage from "@/assets/hero-medical.jpg";

interface HeroSectionProps {
  onGetStarted: () => void;
}

export const HeroSection = ({ onGetStarted }: HeroSectionProps) => {
  return (
    <section className="relative py-20 px-6 overflow-hidden">
      {/* Background gradient */}
      <div className="absolute inset-0 bg-gradient-to-br from-primary-soft via-secondary/30 to-background" />
      
      <div className="relative max-w-7xl mx-auto">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          {/* Content */}
          <div className="space-y-8">
            <div className="space-y-4">
              <h1 className="text-5xl lg:text-6xl font-bold text-foreground leading-tight">
                Understand Your 
                <span className="text-transparent bg-gradient-to-r from-primary to-secondary-accent bg-clip-text">
                  {" "}Medical Reports
                </span>
                <br />
                Simply.
              </h1>
              
              <p className="text-xl text-muted-foreground max-w-2xl leading-relaxed">
                Upload your lab results and get clear, easy-to-understand explanations. 
                No more confusion about what your numbers mean.
              </p>
            </div>
            
            {/* Trust indicators */}
            <div className="flex items-center space-x-6 text-sm text-muted-foreground">
              <div className="flex items-center space-x-2">
                <Shield className="h-4 w-4 text-health-normal" />
                <span>Secure & Private</span>
              </div>
              <div className="flex items-center space-x-2">
                <Zap className="h-4 w-4 text-health-normal" />
                <span>AI-Powered</span>
              </div>
              <div className="flex items-center space-x-2">
                <Heart className="h-4 w-4 text-health-normal" />
                <span>Patient-Friendly</span>
              </div>
            </div>
            
            <Button 
              size="lg" 
              onClick={onGetStarted}
              className="bg-primary hover:bg-primary/90 text-white px-8 py-4 text-lg shadow-floating hover:shadow-soft transition-smooth group"
            >
              Analyze My Report
              <ArrowRight className="ml-2 h-5 w-5 group-hover:translate-x-1 transition-smooth" />
            </Button>
          </div>
          
          {/* Hero image */}
          <div className="relative">
            <div className="relative rounded-3xl overflow-hidden shadow-floating">
              <img 
                src={heroImage} 
                alt="Medical professional with digital health information"
                className="w-full h-auto object-cover"
              />
            </div>
            
            {/* Floating elements */}
            <div className="absolute -top-4 -left-4 bg-white/90 backdrop-blur-sm rounded-2xl p-4 shadow-card">
              <div className="text-sm font-medium text-foreground">Lab Results</div>
              <div className="text-xs text-muted-foreground">Simplified instantly</div>
            </div>
            
            <div className="absolute -bottom-4 -right-4 bg-health-normal-bg/90 backdrop-blur-sm rounded-2xl p-4 shadow-card">
              <div className="text-sm font-medium text-health-normal">All Normal ✓</div>
              <div className="text-xs text-muted-foreground">Easy to understand</div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};