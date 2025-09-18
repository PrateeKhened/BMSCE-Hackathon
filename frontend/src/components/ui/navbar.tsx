import { Button } from "@/components/ui/button";
import { Heart } from "lucide-react";
import { Link } from "react-router-dom";

export const Navbar = () => {
  return (
    <nav className="flex items-center justify-between p-6 bg-background/80 backdrop-blur-sm border-b border-border/50">
      <div className="flex items-center space-x-2">
        <div className="p-2 bg-primary/10 rounded-xl">
          <Heart className="h-6 w-6 text-primary" />
        </div>
        <span className="text-xl font-semibold text-foreground">MedSimple</span>
      </div>
      
      <div className="flex items-center space-x-3">
        <Link to="/login">
          <Button variant="ghost" className="transition-smooth">
            Sign In
          </Button>
        </Link>
        <Link to="/signup">
          <Button className="transition-smooth shadow-soft hover:shadow-card">
            Get Started
          </Button>
        </Link>
      </div>
    </nav>
  );
};
