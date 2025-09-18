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

export function SignupPage() {
  const navigate = useNavigate();
  const { toast } = useToast();
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [dob, setDob] = useState(""); // Date of Birth
  const [gender, setGender] = useState(""); // Gender

  const handleSignup = async () => {
    try {
      const response = await fetch("/api/auth/signup", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          full_name: `${firstName} ${lastName}`,
          email,
          password,
          dob,
          gender,
        }),
      });

      toast({
        title: "Signup successful",
        description: response.message || "You can now log in.",
      });
      navigate("/login");
    } catch (error) {
      if (error instanceof ApiError) {
        toast({
          title: "Signup failed",
          description: error.message,
          variant: "destructive",
        });
      } else {
        toast({
          title: "Signup failed",
          description: "Network error. Please try again.",
          variant: "destructive",
        });
      }
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

      {/* Signup Card */}
      <Card
        className="relative max-w-[400px] w-full backdrop-blur-xl border border-white/30 rounded-3xl shadow-2xl p-8 z-10"
        style={{
          background: "#ffffff", // Solid white background
        }}
      >
        <CardHeader className="text-center mb-8">
          <CardTitle className="text-3xl font-extrabold text-teal-900 drop-shadow-sm">
            Create an Account
          </CardTitle>
          <CardDescription className="text-teal-700 mt-2">
            Fill in your details to get started
          </CardDescription>
        </CardHeader>

        <CardContent className="space-y-6">
          <div className="space-y-4">
            <div className="space-y-1">
              <Label
                htmlFor="first-name"
                className="text-teal-900 font-semibold tracking-wide"
              >
                First Name
              </Label>
              <Input
                id="first-name"
                placeholder="Max"
                className="rounded-lg bg-white/80 text-teal-900 placeholder-teal-600 shadow-inner focus:ring-teal-500 focus:ring-2"
                value={firstName}
                onChange={(e) => setFirstName(e.target.value)}
              />
            </div>

            <div className="space-y-1">
              <Label
                htmlFor="last-name"
                className="text-teal-900 font-semibold tracking-wide"
              >
                Last Name
              </Label>
              <Input
                id="last-name"
                placeholder="Robinson"
                className="rounded-lg bg-white/80 text-teal-900 placeholder-teal-600 shadow-inner focus:ring-teal-500 focus:ring-2"
                value={lastName}
                onChange={(e) => setLastName(e.target.value)}
              />
            </div>

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
              <Label
                htmlFor="password"
                className="text-teal-900 font-semibold tracking-wide"
              >
                Password
              </Label>
              <Input
                id="password"
                type="password"
                placeholder=""
                className="rounded-lg bg-white/80 text-teal-900 placeholder-teal-600 shadow-inner focus:ring-teal-500 focus:ring-2"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>

            {/* New DOB Field */}
            <div className="space-y-1">
              <Label
                htmlFor="dob"
                className="text-teal-900 font-semibold tracking-wide"
              >
                Date of Birth
              </Label>
              <Input
                id="dob"
                type="date"
                className="rounded-lg bg-white/80 text-teal-900 placeholder-teal-600 shadow-inner focus:ring-teal-500 focus:ring-2"
                value={dob}
                onChange={(e) => setDob(e.target.value)}
              />
            </div>

            {/* New Gender Field */}
            <div className="space-y-1">
              <Label
                htmlFor="gender"
                className="text-teal-900 font-semibold tracking-wide"
              >
                Gender
              </Label>
              <select
                id="gender"
                className="w-full rounded-lg bg-white/80 text-teal-900 placeholder-teal-600 shadow-inner focus:ring-teal-500 focus:ring-2"
                value={gender}
                onChange={(e) => setGender(e.target.value)}
              >
                <option value="">Select Gender</option>
                <option value="male">Male</option>
                <option value="female">Female</option>
                <option value="other">Other</option>
              </select>
            </div>
          </div>

          <Button
            onClick={handleSignup}
            className="w-full bg-gradient-to-r from-teal-600 to-cyan-600 text-white font-semibold rounded-lg shadow-md
              hover:from-teal-700 hover:to-cyan-700 transition-all duration-300"
          >
            Create Account
          </Button>

          {/* OR divider */}
          <div className="flex items-center justify-center gap-3 text-teal-700 font-semibold">
            <hr className="flex-grow border-teal-600/30" />
            <span>OR</span>
            <hr className="flex-grow border-teal-600/30" />
          </div>

          {/* Google Sign-Up Button */}
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
            Sign up with Google
          </Button>

          {/* Login Link */}
          <p className="mt-5 text-center text-teal-700">
            Already have an account?{" "}
            <Link
              to="/login"
              className="text-teal-900 font-semibold hover:underline"
            >
              Sign in
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
