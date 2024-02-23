import React from 'react';
import ChessBoard from '../Assets/ChessBoard';
import './MainPage.css';

const MainPage = () => {
    let FEN = 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR';
    //FEN = "3r2k1/7p/2p3p1/1p1qP3/1P2p3/3pP2P/5PP1/Q1R3K1 b - - 5 36";
    FEN = "3k4/8/3b4/3P1N2/2pr4/3QN1B1/8/7K w - - 0 1";
    return (
        <div className="main-page">
            <h1>Chess board</h1>
            <div className='chess-board-wrapper'>
                <ChessBoard fen={FEN} />
            </div>           
        </div>
    )
}

export default MainPage;