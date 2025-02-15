import React, { useEffect, useState } from "react";
import { Navigate } from "react-router-dom";

const ProtectedRoute = ({ children }) => {
    const [isAuthenticated, setIsAuthenticated] = useState(null); // Initially null to show loading state

    useEffect(() => {
        const verifyToken = async () => {
            const token = getCookie("access_token");

            if (!token) {
                setIsAuthenticated(false);
                return;
            }

            try {
                const response = await fetch(
                    `${process.env.REACT_APP_BACKEND_URL}users/verify-token`,
                    {
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json",
                        },
                        body: JSON.stringify({
                            access_token: token,
                            token_type: "bearer",
                        }),
                    }
                );

                if (!response.ok) {
                    throw new Error("Token verification failed");
                }

                setIsAuthenticated(true);
            } catch (error) {
                console.error("Token verification error:", error);
                setIsAuthenticated(false);
            }
        };

        verifyToken();
    }, []);

    const getCookie = (name) => {
        let cookieArr = document.cookie.split(";");
        for (let i = 0; i < cookieArr.length; i++) {
            let cookiePair = cookieArr[i].split("=");
            if (name === cookiePair[0].trim()) {
                return decodeURIComponent(cookiePair[1]);
            }
        }
        return null;
    };

    if (isAuthenticated === null) {
        return <div>Loading...</div>; // Show a loading state while verifying
    }

    if (!isAuthenticated) {
        return <Navigate to="/login" />;
    }

    return children;
};

export default ProtectedRoute;
