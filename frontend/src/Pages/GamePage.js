// import React, { useEffect, useState, useRef } from 'react';
// import { useNavigate, useLocation } from 'react-router-dom';
// import ChessBoard from '../Assets/ChessBoard';
// import './GamePage.css';
// import Header from '../Assets/Header';
// import Footer from '../Assets/Footer';

// const GamePage = () => {
//     const [fen, setFen] = useState('rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1');
//     const [possibleMoves, setPossibleMoves] = useState([]);
//     const [ws, setWs] = useState(null);
//     const [TurnToPlay, setTurnToPlay] = useState('');
//     const [whiteTime, setWhiteTime] = useState(180);
//     const [blackTime, setBlackTime] = useState(180);
//     const [increment, setIncrement] = useState(2);
//     const [openOpenings, setOpenOpenings] = useState([]); // Manage multiple open openings
//     const [selectedVariations, setSelectedVariations] = useState([]); // Allow multiple variations
//     const [openingsData, setOpeningsData] = useState([]); // Initialize as empty array
//     const [showOpenings, setShowOpenings] = useState(false); // Manage openings visibility
//     const timerRef = useRef(null);
//     const navigate = useNavigate();
//     const location = useLocation();

//     useEffect(() => {
//         const params = new URLSearchParams(location.search);
//         const baseTime = parseInt(params.get('baseTime'), 10);
//         const increment = parseInt(params.get('increment'), 10);

//         if (Number.isInteger(baseTime)){
//             setWhiteTime(baseTime);
//             setBlackTime(baseTime);  
//         }
//         if (Number.isInteger(increment)){
//             setIncrement(increment);
//         }
//         // Fetch openings data from the backend
//         fetch(`${process.env.REACT_APP_BACKEND_URL}games/openings`)
//             .then(response => response.json())
//             .then(data => setOpeningsData(data.openings))
//             .catch(error => console.error('Error fetching openings data:', error));

//         // Establish WebSocket connection
//         const websocket = new WebSocket(`ws://localhost:8000/ws/move`);

//         websocket.onopen = () => {
//             console.log('WebSocket connection opened');
//         };

//         websocket.onmessage = (event) => {
//             const data = JSON.parse(event.data);
//             console.log('Message from server:', data);

//             if (data.isUserTurn) {
//                 setTurnToPlay('User');
//                 clearInterval(timerRef.current);
//                 setBlackTime((prevTime) => prevTime + increment); // Add increment after user's move
//                 startTimer('white');
//             } else {
//                 setTurnToPlay('Computer');
//                 clearInterval(timerRef.current);
//                 setWhiteTime((prevTime) => prevTime + increment); // Add increment after computer's move
//                 startTimer('black');
//                 setPossibleMoves(data.possibleMoves);
//             }
//             setFen(data.fen);
//         };

//         websocket.onerror = (event) => {
//             console.error('WebSocket error:', event);
//         };

//         websocket.onclose = () => {
//             console.log('WebSocket is closed now.');
//         };

//         setWs(websocket);

//         return () => {
//             websocket.close();
//             clearInterval(timerRef.current);
//         };
//     }, [location.search]);

//     const startTimer = (side) => {
//         timerRef.current = setInterval(() => {
//             if (side === 'white') {
//                 setWhiteTime((prevTime) => {
//                     if (prevTime <= 0) {
//                         clearInterval(timerRef.current);
//                         return 0;
//                     }
//                     return prevTime - 1;
//                 });
//             } else {
//                 setBlackTime((prevTime) => {
//                     if (prevTime <= 0) {
//                         clearInterval(timerRef.current);
//                         return 0;
//                     }
//                     return prevTime - 1;
//                 });
//             }
//         }, 1000);
//     };

//     const handlePieceDrop = (fromSquare, toSquare) => {
//         console.log("Dropping piece...");
//         const moveData = {
//             fen: fen,
//             from: fromSquare,
//             to: toSquare,
//             selectedVariations: selectedVariations // Send multiple variations
//         };

//         ws.send(JSON.stringify(moveData));
//         console.log("Sending move to server via WebSocket...");
//     };

//     const formatTime = (seconds) => {
//         const minutes = Math.floor(seconds / 60);
//         const secs = seconds % 60;
//         return `${minutes}:${secs < 10 ? '0' : ''}${secs}`;
//     };

//     const toggleOpening = (openingName) => {
//         setOpenOpenings(prev =>
//             prev.includes(openingName) ?
//             prev.filter(name => name !== openingName) :
//             [...prev, openingName]
//         );
//     };

//     const toggleVariation = (variationName) => {
//         setSelectedVariations(prev =>
//             prev.includes(variationName) ?
//             prev.filter(name => name !== variationName) :
//             [...prev, variationName]
//         );
//     };

//     const isOpeningActive = (opening) => {
//         return opening.variations.some(variation => selectedVariations.includes(variation.name));
//     };

//     return (
//         <div>
//             <Header />
//             <section className="game-section">
//                 <div className='game-page-content'>
//                     <div className="clocks-wrapper">
//                         <div className="openings-toggle-wrapper">
//                             <button 
//                                 onClick={() => setShowOpenings(!showOpenings)}
//                                 className={showOpenings ? 'active-opening' : ''}
//                             >
//                                 Select Openings
//                             </button>
//                             <button 
//                                 onClick={() => setShowOpenings(false)}
//                                 className={!showOpenings ? 'active-opening' : ''}
//                             >
//                                 All Openings
//                             </button>
//                         </div>
//                         {showOpenings && (
//                             <div className="openings-wrapper">
//                                 <h3>Openings</h3>
//                                 {openingsData.length > 0 ? (
//                                     openingsData.map((opening, index) => (
//                                         <div key={index}>
//                                             <button
//                                                 onClick={() => toggleOpening(opening.name)}
//                                                 className={isOpeningActive(opening) ? 'active-opening' : ''}
//                                             >
//                                                 {opening.name}
//                                             </button>
//                                             {openOpenings.includes(opening.name) && (
//                                                 <div className="variations-wrapper">
//                                                     {opening.variations.map((variation, vIndex) => (
//                                                         <button
//                                                             key={vIndex}
//                                                             className={selectedVariations.includes(variation.name) ? 'active-variation' : ''}
//                                                             onClick={() => toggleVariation(variation.name)}
//                                                         >
//                                                             {variation.name}
//                                                         </button>
//                                                     ))}
//                                                 </div>
//                                             )}
//                                         </div>
//                                     ))
//                                 ) : (
//                                     <p>Loading openings...</p>
//                                 )}
//                             </div>
//                         )}
//                     </div>
//                     <div className='chess-board-wrapper'>
//                         <div className="player-headings-wrapper">
//                             <p className="player-heading">Chess Engine</p>
//                             <p className="player-subheading" style={{ opacity: TurnToPlay === 'Computer' ? '100%' : '0%' }}>Computer is thinking ...</p>
//                         </div>
//                         <div className='chess-board'>
//                             <ChessBoard fen={fen} onPieceDrop={handlePieceDrop} />
//                         </div>
//                         <div className="player-headings-wrapper">
//                             <p className="player-heading">OKmeray</p>
//                             <p className="player-subheading" style={{ opacity: TurnToPlay === 'User' ? '100%' : '0%' }}>Your move</p>
//                         </div>
//                     </div>
//                     <div className="game-details-wrapper">
//                         <div className={`clock ${TurnToPlay === 'Computer' ? 'active' : ''}`}>
//                             <h2>Black</h2>
//                             <p>{formatTime(blackTime)}</p>
//                         </div>
//                         <div className={`clock ${TurnToPlay === 'User' ? 'active' : ''}`}>
//                             <h2>White</h2>
//                             <p>{formatTime(whiteTime)}</p>
//                         </div>
//                     </div>
//                 </div>
//             </section>
//             <Footer />
//         </div>
//     );
// }

// export default GamePage;

import React, { useEffect, useState, useRef } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import ChessBoard from '../Assets/ChessBoard';
import './GamePage.css';
import Header from '../Assets/Header';
import Footer from '../Assets/Footer';
import Lenis from '@studio-freight/lenis';

const GamePage = () => {
    const [fen, setFen] = useState('rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1');
    const [possibleMoves, setPossibleMoves] = useState([]);
    const [ws, setWs] = useState(null);
    const [TurnToPlay, setTurnToPlay] = useState('');
    const [whiteTime, setWhiteTime] = useState(180);
    const [blackTime, setBlackTime] = useState(180);
    const [increment, setIncrement] = useState(2);
    const [openOpenings, setOpenOpenings] = useState([]); // Manage multiple open openings
    const [selectedVariations, setSelectedVariations] = useState([]); // Allow multiple variations
    const [openingsData, setOpeningsData] = useState([]); // Initialize as empty array
    const [showOpenings, setShowOpenings] = useState(false); // Manage openings visibility
    const timerRef = useRef(null);
    const navigate = useNavigate();
    const location = useLocation();

    useEffect(() => {
        const params = new URLSearchParams(location.search);
        const baseTime = parseInt(params.get('baseTime'), 10);
        const increment = parseInt(params.get('increment'), 10);

        if (Number.isInteger(baseTime)) {
            setWhiteTime(baseTime);
            setBlackTime(baseTime);
        }
        if (Number.isInteger(increment)) {
            setIncrement(increment);
        }

        // Fetch openings data from the backend
        fetch(`${process.env.REACT_APP_BACKEND_URL}games/openings`)
            .then(response => response.json())
            .then(data => setOpeningsData(data.openings))
            .catch(error => console.error('Error fetching openings data:', error));

        // Establish WebSocket connection
        const websocket = new WebSocket(`${process.env.REACT_APP_SOCKET_URL}move`);

        websocket.onopen = () => {
            console.log('WebSocket connection opened');
        };

        websocket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            console.log('Message from server:', data);

            if (data.isUserTurn) {
                setTurnToPlay('User');
                clearInterval(timerRef.current);
                setBlackTime((prevTime) => prevTime + increment); // Add increment after user's move
                startTimer('white');
            } else {
                setTurnToPlay('Computer');
                clearInterval(timerRef.current);
                setWhiteTime((prevTime) => prevTime + increment); // Add increment after computer's move
                startTimer('black');
                setPossibleMoves(data.possibleMoves);
            }
            setFen(data.fen);
        };

        websocket.onerror = (event) => {
            console.error('WebSocket error:', event);
        };

        websocket.onclose = () => {
            console.log('WebSocket is closed now.');
        };

        setWs(websocket);

        return () => {
            websocket.close();
            clearInterval(timerRef.current);
        };
    }, [location.search]);

    useEffect(() => {
        const lenis = new Lenis({
            duration: 1.2,
            easing: (t) => Math.min(1, 1.001 - Math.pow(2, -10 * t)),
        });

        function raf(time) {
            lenis.raf(time);
            requestAnimationFrame(raf);
        }

        requestAnimationFrame(raf);
        return () => lenis.destroy();
    }, []);

    const startTimer = (side) => {
        timerRef.current = setInterval(() => {
            if (side === 'white') {
                setWhiteTime((prevTime) => {
                    if (prevTime <= 0) {
                        clearInterval(timerRef.current);
                        return 0;
                    }
                    return prevTime - 1;
                });
            } else {
                setBlackTime((prevTime) => {
                    if (prevTime <= 0) {
                        clearInterval(timerRef.current);
                        return 0;
                    }
                    return prevTime - 1;
                });
            }
        }, 1000);
    };

    const handlePieceDrop = (fromSquare, toSquare) => {
        console.log("Dropping piece...");
        const moveData = {
            fen: fen,
            from: fromSquare,
            to: toSquare,
            selectedVariations: selectedVariations // Send multiple variations
        };

        ws.send(JSON.stringify(moveData));
        console.log("Sending move to server via WebSocket...");
    };

    const formatTime = (seconds) => {
        const minutes = Math.floor(seconds / 60);
        const secs = seconds % 60;
        return `${minutes}:${secs < 10 ? '0' : ''}${secs}`;
    };

    const toggleOpening = (openingName) => {
        setOpenOpenings(prev =>
            prev.includes(openingName) ?
            prev.filter(name => name !== openingName) :
            [...prev, openingName]
        );
    };

    const toggleVariation = (variationName) => {
        setSelectedVariations(prev =>
            prev.includes(variationName) ?
            prev.filter(name => name !== variationName) :
            [...prev, variationName]
        );
    };

    const isOpeningActive = (opening) => {
        return opening.variations.some(variation => selectedVariations.includes(variation.name));
    };

    return (
        <div>
            <Header />
            <section className="game-section">
                <div className='game-page-content'>
                    <div className="clocks-wrapper">
                        <div className="openings-toggle-wrapper">
                            <button
                                onClick={() => setShowOpenings(!showOpenings)}
                                className={showOpenings ? 'active-opening' : ''}
                            >
                                Select Openings
                            </button>
                            <button
                                onClick={() => setShowOpenings(false)}
                                className={!showOpenings ? 'active-opening' : ''}
                            >
                                All Openings
                            </button>
                        </div>
                        {showOpenings && (
                            <div className="openings-wrapper">
                                <h3>Openings</h3>
                                {openingsData.length > 0 ? (
                                    openingsData.map((opening, index) => (
                                        <div key={index}>
                                            <button
                                                onClick={() => toggleOpening(opening.name)}
                                                className={isOpeningActive(opening) ? 'active-opening' : ''}
                                            >
                                                {opening.name}
                                            </button>
                                            {openOpenings.includes(opening.name) && (
                                                <div className="variations-wrapper">
                                                    {opening.variations.map((variation, vIndex) => (
                                                        <button
                                                            key={vIndex}
                                                            className={selectedVariations.includes(variation.name) ? 'active-variation' : ''}
                                                            onClick={() => toggleVariation(variation.name)}
                                                        >
                                                            {variation.name}
                                                        </button>
                                                    ))}
                                                </div>
                                            )}
                                        </div>
                                    ))
                                ) : (
                                    <p>Loading openings...</p>
                                )}
                            </div>
                        )}
                    </div>
                    <div className='chess-board-wrapper'>
                        <div className="player-headings-wrapper">
                            <p className="player-heading">Chess Engine</p>
                            <p className="player-subheading" style={{ opacity: TurnToPlay === 'Computer' ? '100%' : '0%' }}>Computer is thinking ...</p>
                        </div>
                        <div className='chess-board'>
                            <ChessBoard fen={fen} onPieceDrop={handlePieceDrop} />
                        </div>
                        <div className="player-headings-wrapper">
                            <p className="player-heading">OKmeray</p>
                            <p className="player-subheading" style={{ opacity: TurnToPlay === 'User' ? '100%' : '0%' }}>Your move</p>
                        </div>
                    </div>
                    <div className="game-details-wrapper">
                        <div className={`clock ${TurnToPlay === 'Computer' ? 'active' : ''}`}>
                            <h2>Black</h2>
                            <p>{formatTime(blackTime)}</p>
                        </div>
                        <div className={`clock ${TurnToPlay === 'User' ? 'active' : ''}`}>
                            <h2>White</h2>
                            <p>{formatTime(whiteTime)}</p>
                        </div>
                    </div>
                </div>
            </section>
            <Footer />
        </div>
    );
}

export default GamePage;
