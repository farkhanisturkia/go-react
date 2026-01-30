import React, { createContext, useState, useEffect, ReactNode } from 'react';
import Cookies from 'js-cookie';

// Menentukan tipe dari context value
interface AuthContextType {
    isAuthenticated: boolean;
    setIsAuthenticated: React.Dispatch<React.SetStateAction<boolean>>;
}

// Membuat context dengan default value undefined
export const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Menentukan tipe props untuk AuthProvider
interface AuthProviderProps {
    children: ReactNode;
}

// Komponen provider untuk konteks otentikasi
export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
    const [isAuthenticated, setIsAuthenticated] = useState<boolean>(!!Cookies.get('token'));

    useEffect(() => {
        const handleTokenChange = () => {
            setIsAuthenticated(!!Cookies.get('token'));
        };

        window.addEventListener('storage', handleTokenChange);
        return () => {
            window.removeEventListener('storage', handleTokenChange);
        };
    }, []);

    return (
        <AuthContext.Provider value={{ isAuthenticated, setIsAuthenticated }}>
            {children}
        </AuthContext.Provider>
    );
};
