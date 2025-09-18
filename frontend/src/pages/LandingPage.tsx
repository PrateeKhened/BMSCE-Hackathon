import { Navbar } from "@/components/ui/navbar";
import { HeroSection } from "@/components/ui/hero-section";
import { useNavigate } from "react-router-dom";

const LandingPage = () => {
  const navigate = useNavigate();

  const handleGetStarted = () => {
    navigate("/analyze");
  };

  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      <main>
        <HeroSection onGetStarted={handleGetStarted} />
        
        {/* Features preview */}
        <section className="py-16 px-6 bg-primary-soft/50">
          <div className="max-w-6xl mx-auto text-center">
            <h2 className="text-3xl font-semibold text-foreground mb-4">
              How It Works
            </h2>
            <p className="text-muted-foreground mb-12 text-lg">
              Three simple steps to understand your health better
            </p>
            
            <div className="grid md:grid-cols-3 gap-8">
              <div className="p-6 bg-card rounded-2xl shadow-card transition-smooth hover:shadow-floating">
                <div className="w-12 h-12 bg-primary/10 rounded-xl flex items-center justify-center mx-auto mb-4">
                  <span className="text-2xl">📋</span>
                </div>
                <h3 className="font-semibold text-lg mb-2">Upload Report</h3>
                <p className="text-muted-foreground">Simply upload your PDF lab results or medical report</p>
              </div>
              
              <div className="p-6 bg-card rounded-2xl shadow-card transition-smooth hover:shadow-floating">
                <div className="w-12 h-12 bg-primary/10 rounded-xl flex items-center justify-center mx-auto mb-4">
                  <span className="text-2xl">🤖</span>
                </div>
                <h3 className="font-semibold text-lg mb-2">AI Analysis</h3>
                <p className="text-muted-foreground">Our AI breaks down complex medical jargon into simple language</p>
              </div>
              
              <div className="p-6 bg-card rounded-2xl shadow-card transition-smooth hover:shadow-floating">
                <div className="w-12 h-12 bg-primary/10 rounded-xl flex items-center justify-center mx-auto mb-4">
                  <span className="text-2xl">💡</span>
                </div>
                <h3 className="font-semibold text-lg mb-2">Clear Insights</h3>
                <p className="text-muted-foreground">Get easy-to-understand explanations of your health metrics</p>
              </div>
            </div>
          </div>
        </section>
      </main>
    </div>
  );
};

export default LandingPage;