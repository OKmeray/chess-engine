import React, { useState } from "react";
import Header from "../Components/Header";
import Footer from "../Components/Footer";
import "./LoginPage.css";

function Signup() {
    const [formData, setFormData] = useState({
        username: "",
        email: "",
        password: "",
    });
    const [errorMessage, setErrorMessage] = useState("");
    const [showError, setShowError] = useState(false);

    const handleChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    const setCookie = (name, value, days) => {
        let expires = "";
        if (days) {
            let date = new Date();
            date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
            expires = "; expires=" + date.toUTCString();
        }
        document.cookie = name + "=" + (value || "") + expires + "; path=/";
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await fetch(
                `${process.env.REACT_APP_BACKEND_URL}users/register`,
                {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({
                        username: formData.username,
                        email: formData.email,
                        password: formData.password,
                    }),
                }
            );

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            } else {
                const result = await response.json();
                setCookie("access_token", result.access_token, 1); // Set cookie to expire in 1 day
                setShowError(false);
            }
        } catch (error) {
            console.error("Failed to register:", error);
            setErrorMessage("Registration failed. Please check your inputs.");
            setShowError(true);
            setFormData({ ...formData, password: "" }); // Clear the password field
        }
    };

    return (
        <div>
            <Header />
            <div className="login-container">
                <div className="login-form">
                    <h2>Sign Up</h2>
                    <div className="welcome-message">
                        Join us! Please fill in the details to create an
                        account.
                    </div>
                    {showError && (
                        <div className="error-message">{errorMessage}</div>
                    )}
                    <form onSubmit={handleSubmit}>
                        <div className="form-group">
                            <label>Username:</label>
                            <input
                                type="text"
                                className="form-control"
                                name="username"
                                value={formData.username}
                                onChange={handleChange}
                            />
                        </div>
                        <div className="form-group">
                            <label>Email:</label>
                            <input
                                type="email"
                                className="form-control"
                                name="email"
                                value={formData.email}
                                onChange={handleChange}
                            />
                        </div>
                        <div className="form-group">
                            <label>Password:</label>
                            <input
                                type="password"
                                className="form-control"
                                name="password"
                                value={formData.password}
                                onChange={handleChange}
                            />
                        </div>
                        <button type="submit" className="btn btn-primary">
                            Sign Up
                        </button>
                    </form>
                </div>
            </div>
            <Footer />
        </div>
    );
}

export default Signup;
