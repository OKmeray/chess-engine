import { SET_CAN_DRAG } from './types';

const initialState = {
  isUserTurn: true
};

const appReducer = (state = initialState, action) => {
  switch (action.type) {
    case SET_CAN_DRAG:
      return {
        ...state,
        isUserTurn: action.payload
      };
    default:
      return state;
  }
};

export default appReducer;