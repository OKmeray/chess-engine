// import React from 'react';
// import './Footer.css';

// const Footer = () => {
//   return (
//     <footer className="footer">
//       <div className="container">
//         <div className="row">
//           <div className="col-md-6">
//             <h5>Contact Us</h5>
//             <p>Email: example@example.com</p>
//             <p>Phone: 123-456-7890</p>
//             <p>Address: 123 Street, City, Country</p>
//           </div>
//           <div className="col-md-6">
//             <h5>Follow Us</h5>
//             <ul className="list-unstyled">
//               <li><a href="https://twitter.com/example" target="_blank" rel="noopener noreferrer">Twitter</a></li>
//               <li><a href="https://facebook.com/example" target="_blank" rel="noopener noreferrer">Facebook</a></li>
//               <li><a href="https://instagram.com/example" target="_blank" rel="noopener noreferrer">Instagram</a></li>
//             </ul>
//           </div>
//         </div>
//       </div>
//     </footer>
//   );
// };

// export default Footer;

import React from 'react';
import './Footer.css';

const Footer = () => {
  return (
    <footer className="footer">
      <div className="footer-container">
        <div className="footer-content">
          <div className="footer-section contact">
            <h5>Contact Us</h5>
            <p>Email: <a href="mailto:chessengine@gmail.com">chessengine@gmail.com</a></p>
            <p>Phone: <a href="tel:+1234567890">123-456-7890</a></p>
            <p>Address: 123 Street, City, Country</p>
          </div>
          <div className="footer-section follow-us">
            <h5>Follow Us</h5>
            <ul className="social-links">
              <li><a href="https://twitter.com/example" target="_blank" rel="noopener noreferrer">Twitter</a></li>
              <li><a href="https://facebook.com/example" target="_blank" rel="noopener noreferrer">Facebook</a></li>
              <li><a href="https://instagram.com/example" target="_blank" rel="noopener noreferrer">Instagram</a></li>
            </ul>
          </div>
          <div className="footer-section about">
            <h5>About Us</h5>
            <p>I am a chess engine developer dedicated to providing top-notch chess software solutions for players of all levels. My mission is to enhance the chess-playing experience with innovative and user-friendly tools.</p>
          </div>
        </div>
        <div className="footer-bottom">
          <p>&copy; {new Date().getFullYear()} Chess Engine. All rights reserved.</p>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
