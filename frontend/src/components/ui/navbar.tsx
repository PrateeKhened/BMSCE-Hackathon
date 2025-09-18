import { Button } from "@/components/ui/button";
import { Heart, FileText } from "lucide-react";

interface NavbarProps {
  onSignIn: () => void;
  onSignUp: () => void;
}

export const Navbar = ({ onSignIn, onSignUp }: NavbarProps) => {
  return (
    <nav className="flex items-center justify-between p-6 bg-background/80 backdrop-blur-sm border-b border-border/50">
      <div className="flex items-center space-x-2">
        <div className="p-2 bg-primary/10 rounded-xl">
          <Heart className="h-6 w-6 text-primary" />
        </div>
        <span className="text-xl font-semibold text-foreground">MedSimple</span>
      </div>
      
      <div className="flex items-center space-x-3">
        <Button variant="ghost" onClick={onSignIn} className="transition-smooth">
          Sign In
        </Button>
        <Button onClick={onSignUp} className="transition-smooth shadow-soft hover:shadow-card">
          Get Started
        </Button>
      </div>
    </nav>
  );
};