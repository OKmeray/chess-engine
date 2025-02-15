import React from "react";
import { Link } from "react-router-dom";
import "./GuestPage.css";
import Header from "../Components/Header";
import Footer from "../Components/Footer";

const GuestPage = () => {
    return (
        <div className="page-container">
            <Header />
            <div className="content-container">
                <div className="hero">
                    <h1>Welcome to Our Chess Engine</h1>
                    <div className="hero-buttons">
                        <Link
                            to="/login"
                            className="btn btn-primary auth-buttons"
                        >
                            Login
                        </Link>
                        <Link
                            to="/signup"
                            className="btn btn-secondary auth-buttons"
                        >
                            Sign Up
                        </Link>
                    </div>
                </div>
                <div className="features">
                    <h2>Features</h2>
                    <div className="feature-list">
                        <div className="feature-item">
                            <h3>Advanced AI</h3>
                            <p>
                                Play against a powerful chess engine that adapts
                                to your skill level.
                            </p>
                        </div>
                        <div className="feature-item">
                            <h3>Opening Theory</h3>
                            <p>
                                Learn and practice opening theory by selecting
                                and training various chess openings directly on
                                the game page.
                            </p>
                        </div>
                        <div className="feature-item">
                            <h3>Interactive Tutorials</h3>
                            <p>
                                Learn and improve your chess skills with our
                                interactive lessons.
                            </p>
                        </div>
                    </div>
                </div>
                <div className="testimonials">
                    <h2>What Our Users Say</h2>
                    <div className="testimonial-list">
                        <div className="testimonial-item">
                            <p>
                                "This chess engine is amazing! It has really
                                helped me improve my game."
                            </p>
                            <h4>- John Doe</h4>
                        </div>
                        <div className="testimonial-item">
                            <p>
                                "I love the interactive tutorials. They are so
                                helpful for beginners like me."
                            </p>
                            <h4>- Jane Smith</h4>
                        </div>
                    </div>
                </div>
            </div>
            <Footer />
        </div>
    );
};

export default GuestPage;
