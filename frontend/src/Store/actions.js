// src/store/actions.js
import { SET_CAN_DRAG } from './types';

export const setCanDrag = (canDrag) => ({
  type: SET_CAN_DRAG,
  payload: canDrag
});
