import React, { useState } from 'react';
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';

import blackKing from '../Assets/Images/king_black.svg';
import blackQueen from '../Assets/Images/queen_black.svg';
import blackRook from '../Assets/Images/rook_black.svg';
import blackBishop from '../Assets/Images/bishop_black.svg';
import blackKnight from '../Assets/Images/knight_black.svg';
import blackPawn from '../Assets/Images/pawn_black.svg';
import whiteKing from '../Assets/Images/king_white.svg';
import whiteQueen from '../Assets/Images/queen_white.svg';
import whiteRook from '../Assets/Images/rook_white.svg';
import whiteBishop from '../Assets/Images/bishop_white.svg';
import whiteKnight from '../Assets/Images/knight_white.svg';
import whitePawn from '../Assets/Images/pawn_white.svg';

import Piece from '../Assets/Piece';
import squareSize from '../Assets/variables';
import Pieces from '../Assets/Pieces';

const ChessBoard = ({fen}) => {

    const boardWrapperStyles = {
        width: squareSize * 8 + "px",
        height: squareSize * 8 + "px",
        position: "relative",
        align: "center",
        display: "grid",
        gridTemplateColumns: 'repeat(8, 1fr)',
        gridTemplateRows: 'repeat(8, 1fr)',
    }

    //fen = "3k4/8/5q2/8/2Q5/4P3/8/3K4 w - - 0 1";
    //     /* get position of the pieces and set it in array */
    const piecePlacement = fen.split(' ')[0];

    function isDigit(char) {
        if (char >= '0' && char <= '9') {
            return true;
        }
        return false;
    }

    let pieces = [];
    let squareIndex = 0;
    let indexCounter = 0;
    for (let i = 0; i < piecePlacement.length; i++) {

        if (isDigit(piecePlacement[i]))
        {
            squareIndex += parseInt(piecePlacement[i]);
        }
        else if (piecePlacement[i] === "/") {
            continue;
        }
        else {
            pieces.push({
                "id": indexCounter,
                "piece": piecePlacement[i],
                "offset": squareIndex
            });
            squareIndex += 1;
            indexCounter += 1;
        }
        
    }

    return (
    <div className='board-image-wrapper' style={boardWrapperStyles}>
        {Array.from({ length: 8 }, (_, row) => (
            Array.from({ length: 8 }, (_, col) => (
            <div style={
                {
                    "background": (row + col) % 2 === 0 ? '#f0d9b5' : '#b58863', 
                    "width": `${squareSize}px`,
                    "height": `${squareSize}px`
                }
            }>
            </div>
            ))
        ))}
        <Pieces pieces={pieces}/>
    </div>
    );
}
export default ChessBoard;
