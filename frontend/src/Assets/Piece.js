import React from 'react';
import squareSize from '../Assets/variables';

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


const Piece = ({piece, square}) => {
  function getImageFromFenSymbol(char) {
    switch (char) {
        case 'k':
            return blackKing;
        case 'q':
            return blackQueen;
        case 'r':
            return blackRook;
        case 'b':
            return blackBishop;
        case 'n':
            return blackKnight;
        case 'p':
            return blackPawn;
        case 'K':
            return whiteKing;
        case 'Q':
            return whiteQueen;
        case 'R':
            return whiteRook;
        case 'B':
            return whiteBishop;
        case 'N':
            return whiteKnight;
        case 'P':
            return whitePawn
    }
  }

  function getImageStyles(piece, square) {
    return {
      position: "absolute",
      top: `${Math.floor(square / 8) * squareSize}px`,
      left: `${(square % 8) * squareSize}px`
    }
  }

  const onDragStart = e => {
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('text/plain', `${piece},${square}`);
    setTimeout(() => {
      e.target.style.display = "none";
    }, 0);
  }

  const onDragOver = e => {
    e.preventDefault();
  }

  const onDragEnd = e => {
    e.target.style.display = "";
  };

  return (
    <img 
      src={getImageFromFenSymbol(piece)}
      width={squareSize}
      height={squareSize}
      style={getImageStyles(piece, square)}
      draggable={true}
      onDragStart={(e) => onDragStart(e)}
      onDragOver={onDragOver}
      onDragEnd={onDragEnd}
    />
  )
};



export default Piece;
