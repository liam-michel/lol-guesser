import { createContext, useContext, useState, useEffect, ReactNode } from "react";
import { loginUser, createUser } from "../utils/user_api";

interface AuthenticationType {
  username: string | null;
  token: string | null;
  login: (username: string, password: string) => Promise<any>;
  createAccount: (username: string, password: string) => Promise<any>;
  logout: () => void;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthenticationType>({} as AuthenticationType);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [username, setUsername] = useState<string | null>(() => localStorage.getItem("authUsername"));
  const [token, setToken] = useState<string | null>(() => localStorage.getItem("authToken"));

  // Sync username and token with localStorage on change
  useEffect(() => {
    if (token) {
      localStorage.setItem("authToken", token);
    }
  }, [token]);

  useEffect(() => {
    if (username) {
      localStorage.setItem("authUsername", username);
    }
  }, [username]);

  const login = async (username: string, password: string) => {
    const response = await loginUser(username, password);
    if (response.error) {
      // Return a consistent error shape with a string error message
      return { error: response.error };
    }
    console.log("Logged in");
    setToken(response.token);
    setUsername(response.username);
    return;
  };

  const createAccount = async (username: string, password: string) => {
    const response = await createUser(username, password);
  
    if (response.error) {
      // Return a consistent error shape with a string error message
      return { error: response.error };
    }
  
    console.log("Created new account");
    setToken(response.token);
    setUsername(response.username);
    return;
    //return { token: response.token, username: response.username }; // Return successful data
  };
  

  const logout = () => {
    setToken(null);
    setUsername(null);
    localStorage.removeItem("authToken");
    localStorage.removeItem("authUsername");
  };

  const isAuthenticated = !!username;

  const value = {
    username,
    token,
    createAccount,
    login,
    logout,
    isAuthenticated
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
