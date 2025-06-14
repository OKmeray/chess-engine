import React from "react";
import "./Footer.css";

const Footer = () => {
    return (
        <footer className="footer">
            <div className="footer-container">
                <div className="footer-content">
                    <div className="footer-section about">
                        <h5>Про розробника</h5>
                        <p>
                            Я студент 4-го курсу Львівського Національного
                            Університету імені Івана Франка факультету
                            прикладної математики та інформатики. Займаюсь
                            шахами більше 10 років і вирішив спробувати написати
                            свій власний шаховий рушій, з яким можу грати та
                            тренуватися.
                        </p>
                    </div>
                    <div className="footer-section contact">
                        <h5>Контакти</h5>
                        <p>
                            Email:{" "}
                            <a href="mailto:roman.yaremko.pmi@lnu.edu.ua">
                                roman.yaremko.pmi@lnu.edu.ua
                            </a>
                        </p>
                        <p>
                            Адреса:{" "}
                            <a
                                href="https://www.google.com/maps/place/Universytetska+St,+1,+L'viv,+L'vivs'ka+oblast,+79000/@49.840224,24.0221738,17z"
                                target="_blank"
                            >
                                Україна, м. Львів, вул. Університетська, 1
                            </a>
                        </p>
                    </div>
                    <div className="footer-section follow-us">
                        <h5>Соціальні мережі</h5>
                        <ul className="social-links">
                            <li>
                                <a
                                    href="https://facebook.com"
                                    target="_blank"
                                    rel="noopener noreferrer"
                                >
                                    Facebook
                                </a>
                            </li>
                            <li>
                                <a
                                    href="https://instagram.com"
                                    target="_blank"
                                    rel="noopener noreferrer"
                                >
                                    Instagram
                                </a>
                            </li>
                        </ul>
                    </div>
                </div>
                <div className="footer-bottom">
                    <p>
                        &copy; {new Date().getFullYear()} Шаховий рушій. Всі
                        права захищено.
                    </p>
                </div>
            </div>
        </footer>
    );
};

export default Footer;
