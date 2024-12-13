import React, { createContext, useContext, useState, useEffect, ReactNode } from "react";


interface AuthenticationType{
  username: {username: string} | null;
  token: string | null;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean

}


const AuthContext = createContext<AuthenticationType>({} as AuthenticationType);

export const AuthProvider = ({ children }: {children: ReactNode}) => {
  
  
}