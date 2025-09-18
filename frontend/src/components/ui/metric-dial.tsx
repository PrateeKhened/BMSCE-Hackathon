import { cn } from "@/lib/utils";

interface MetricDialProps {
  title: string;
  value: string;
  normalRange: string;
  percentage: number;
  status: "normal" | "warning" | "danger";
  className?: string;
}

export const MetricDial = ({ 
  title, 
  value, 
  normalRange, 
  percentage, 
  status, 
  className 
}: MetricDialProps) => {
  const getStatusColors = (status: string) => {
    switch (status) {
      case "normal":
        return {
          stroke: "stroke-health-normal",
          bg: "bg-health-normal-bg",
          text: "text-health-normal"
        };
      case "warning":
        return {
          stroke: "stroke-health-warning",
          bg: "bg-health-warning-bg",
          text: "text-health-warning"
        };
      case "danger":
        return {
          stroke: "stroke-health-danger",
          bg: "bg-health-danger-bg",
          text: "text-health-danger"
        };
      default:
        return {
          stroke: "stroke-muted",
          bg: "bg-muted",
          text: "text-muted-foreground"
        };
    }
  };

  const colors = getStatusColors(status);
  const circumference = 2 * Math.PI * 40; // radius = 40
  const strokeDasharray = circumference;
  const strokeDashoffset = circumference - (percentage / 100) * circumference;

  return (
    <div className={cn("flex flex-col items-center p-6 rounded-2xl transition-smooth hover:shadow-card", colors.bg, className)}>
      <div className="relative w-28 h-28 mb-4">
        <svg className="transform -rotate-90 w-28 h-28" viewBox="0 0 100 100">
          {/* Background circle */}
          <circle
            cx="50"
            cy="50"
            r="40"
            stroke="currentColor"
            strokeWidth="6"
            fill="none"
            className="text-border opacity-30"
          />
          {/* Progress circle */}
          <circle
            cx="50"
            cy="50"
            r="40"
            stroke="currentColor"
            strokeWidth="6"
            fill="none"
            strokeLinecap="round"
            strokeDasharray={strokeDasharray}
            strokeDashoffset={strokeDashoffset}
            className={cn("transition-all duration-1000 ease-out", colors.stroke)}
            style={{
              animationDelay: '0.5s'
            }}
          />
        </svg>
        
        {/* Center content */}
        <div className="absolute inset-0 flex flex-col items-center justify-center">
          <span className={cn("text-xl font-bold", colors.text)}>{value}</span>
          <span className="text-xs text-muted-foreground">{percentage}%</span>
        </div>
      </div>
      
      <div className="text-center space-y-1">
        <h3 className="font-semibold text-foreground">{title}</h3>
        <p className="text-xs text-muted-foreground">Normal: {normalRange}</p>
      </div>
    </div>
  );
};