import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Link, useNavigate } from "react-router-dom";
import { useToast } from "@/hooks/use-toast";
import { authApi, ApiError } from "@/lib/api";

export function LoginPage() {
  const navigate = useNavigate();
  const { toast } = useToast();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const handleLogin = async () => {
    setLoading(true);
    try {
      const response = await authApi.login({ email, password });

      if (response.token) {
        authApi.setToken(response.token);
        toast({
          title: "Login successful",
          description: response.message || "Welcome back!",
        });
        navigate("/analyze");
      } else {
        toast({
          title: "Login failed",
          description: errorData.error || "Invalid credentials.",
          variant: "destructive",
        });
      }
    } catch (error) {
      toast({
        title: "Error",
        description: "Something went wrong. Please try again.",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      className="relative min-h-screen flex items-center justify-center px-6"
      style={{
        background: "var(--gradient-hero)",
        backgroundRepeat: "no-repeat",
        backgroundAttachment: "fixed",
      }}
    >
      {/* Soft abstract shapes */}
      <svg
        className="absolute top-0 left-0 w-full h-full pointer-events-none"
        xmlns="http://www.w3.org/2000/svg"
        preserveAspectRatio="xMidYMid meet"
        fill="none"
      >
        <circle cx="20%" cy="15%" r="180" fill="rgba(255,255,255,0.06)" />
        <circle cx="80%" cy="80%" r="260" fill="rgba(255,255,255,0.08)" />
        <circle cx="50%" cy="50%" r="350" fill="rgba(255,255,255,0.03)" />
      </svg>

      {/* Glassy translucent white card */}
      <Card
        className="relative max-w-[400px] w-full backdrop-blur-xl border border-white/30 rounded-3xl shadow-2xl p-8 z-10"
        style={{
          background: "#ffffff", // translucent white gradient
          // You can replace above with var(--gradient-card) if you want color instead of white gradient
          // background: "var(--gradient-card)",
        }}
      >
        <CardHeader className="text-center mb-8">
          <CardTitle className="text-3xl font-extrabold text-teal-900 drop-shadow-sm">
            Welcome Back
          </CardTitle>
          <CardDescription className="text-teal-700 mt-2">
            Login to your account and start exploring
          </CardDescription>
        </CardHeader>

        <CardContent className="space-y-6">
          <div className="space-y-4">
            <div className="space-y-1">
              <Label
                htmlFor="email"
                className="text-teal-900 font-semibold tracking-wide"
              >
                Email Address
              </Label>
              <Input
                id="email"
                type="email"
                placeholder="you@example.com"
                className="rounded-lg bg-white/80 text-teal-900 placeholder-teal-600 shadow-inner focus:ring-teal-500 focus:ring-2"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>

            <div className="space-y-1">
              <div className="flex justify-between items-center">
                <Label
                  htmlFor="password"
                  className="text-teal-900 font-semibold tracking-wide"
                >
                  Password
                </Label>
                <Link
                  to="#"
                  className="text-teal-700 hover:text-teal-900 text-sm font-medium"
                >
                  Forgot Password?
                </Link>
              </div>
              <Input
                id="password"
                type="password"
                placeholder=""
                className="rounded-lg bg-white/80 text-teal-900 placeholder-teal-600 shadow-inner focus:ring-teal-500 focus:ring-2"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
          </div>

          <Button
            onClick={handleLogin}
            disabled={loading}
            className="w-full bg-gradient-to-r from-teal-600 to-cyan-600 text-white font-semibold rounded-lg shadow-md
              hover:from-teal-700 hover:to-cyan-700 transition-all duration-300"
          >
            {loading ? "Signing In..." : "Sign In"}
          </Button>

          {/* OR divider */}
          <div className="flex items-center justify-center gap-3 text-teal-700 font-semibold">
            <hr className="flex-grow border-teal-600/30" />
            <span>OR</span>
            <hr className="flex-grow border-teal-600/30" />
          </div>

          {/* Google Login Button */}
          <Button
            variant="outline"
            className="w-full flex items-center justify-center gap-3 border-teal-600 text-teal-700
              hover:bg-teal-600 hover:text-white transition-colors rounded-lg"
          >
            <img
              src="https://www.svgrepo.com/show/475656/google-color.svg"
              alt="Google"
              className="w-6 h-6"
            />
            Continue with Google
          </Button>

          {/* Signup Link */}
          <p className="mt-5 text-center text-teal-700">
            Don&apos;t have an account?{" "}
            <Link
              to="/signup"
              className="text-teal-900 font-semibold hover:underline"
            >
              Sign up
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}