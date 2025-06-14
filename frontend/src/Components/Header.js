import React from "react";
import { Link } from "react-router-dom";
import chessLogo from "../assets/images/king_white.svg";
import "./Header.css";

const Header = () => {
    return (
        <nav className="navbar">
            <div className="container">
                <Link className="navbar-brand" to="/">
                    <img
                        src={chessLogo}
                        width="30"
                        height="30"
                        className="d-inline-block align-top mr-2"
                        alt="Chess Board Logo"
                    />
                    Шаховий рушій
                </Link>
                <nav className="nav-links">
                    <Link className="nav-link" to="/">
                        Головна
                    </Link>
                    <Link className="nav-link" to="/game">
                        Грати
                    </Link>
                    <Link className="nav-link" to="/models">
                        Моделі
                    </Link>
                </nav>
                <Link className="header-button" to="/models">
                    <div className="header-button-text">Моделі нейромереж</div>
                    <div className="button-circle">
                        <svg
                            width="14"
                            height="14"
                            viewBox="0 0 14 14"
                            fill="none"
                            xmlns="http://www.w3.org/2000/svg"
                        >
                            <path
                                d="M13.5546 1.13689L13.5546 10.6828C13.5546 10.8819 13.4756 11.0727 13.3349 11.2135C13.1941 11.3542 13.0032 11.4333 12.8042 11.4333C12.6052 11.4333 12.4143 11.3542 12.2736 11.2135C12.1329 11.0727 12.0538 10.8819 12.0538 10.6828L12.0545 2.94731L1.66729 13.3345C1.52664 13.4751 1.33588 13.5542 1.13696 13.5542C0.938051 13.5542 0.747285 13.4751 0.606632 13.3345C0.46598 13.1938 0.386963 13.0031 0.386963 12.8042C0.386963 12.6052 0.46598 12.4145 0.606632 12.2738L10.9938 1.88665L3.25828 1.88731C3.05926 1.88731 2.86839 1.80825 2.72766 1.66752C2.58693 1.52679 2.50787 1.33592 2.50787 1.13689C2.50787 0.937869 2.58693 0.746998 2.72766 0.606267C2.86839 0.465536 3.05926 0.386475 3.25828 0.386475L12.8042 0.386475C12.9028 0.38642 13.0004 0.405792 13.0915 0.443485C13.1825 0.481178 13.2653 0.536452 13.335 0.606145C13.4047 0.675839 13.4599 0.758585 13.4976 0.849655C13.5353 0.940723 13.5547 1.03833 13.5546 1.13689Z"
                                fill="#fff"
                            />
                        </svg>
                    </div>
                </Link>
            </div>
        </nav>
    );
};

export default Header;
