import React from "react";
import { Link, useNavigate } from "react-router-dom";
import chessLogo from "../Assets/Images/king_white.svg";
import "./Header.css";

const Header = () => {
    const navigate = useNavigate();

    const handleLogout = () => {
        // Remove the access_token from cookies
        document.cookie = "access_token=; Max-Age=0; path=/;";

        navigate("/");
    };

    const isLoggedIn = () => {
        return document.cookie
            .split(";")
            .some((item) => item.trim().startsWith("access_token="));
    };

    return (
        <nav className="navbar">
            <div className="container">
                <a className="navbar-brand" href="/home">
                    <img
                        src={chessLogo}
                        width="30"
                        height="30"
                        className="d-inline-block align-top mr-2"
                        alt="Chess Board Logo"
                    />
                    Chess Engine
                </a>

                <ul className="navbar-nav ml-auto">
                    {isLoggedIn() ? (
                        <>
                            <li className="nav-item">
                                <Link className="nav-link" to="/home">
                                    Home
                                </Link>
                            </li>
                            <li className="nav-item">
                                <Link
                                    className="nav-link"
                                    to="/game?baseTime=180&increment=2"
                                >
                                    Play
                                </Link>
                            </li>
                            {/* <li className="nav-item">
                <Link className="nav-link" to="/profile">Profile</Link>
              </li> */}
                            <li className="nav-item">
                                <button
                                    className="nav-link btn"
                                    onClick={handleLogout}
                                >
                                    Logout
                                </button>
                            </li>
                        </>
                    ) : (
                        <>
                            <li className="nav-item">
                                <Link className="nav-link" to="/login">
                                    Login
                                </Link>
                            </li>
                            <li className="nav-item">
                                <Link className="nav-link" to="/signup">
                                    Sign Up
                                </Link>
                            </li>
                        </>
                    )}
                </ul>
            </div>
        </nav>
    );
};

export default Header;
