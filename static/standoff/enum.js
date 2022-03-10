import { globalEn } from '../globalEnum.js';

var en = {
    Phase: {
        PREGAME: '0',
        PLAY: '1',
    },
    ToServerCode: {
        START_GAME: '0',
        END_GAME: '2',
        DECISION: 'a',
        PROMPT_REQUEST: 'b',
    },
    ToClientCode: {
        START_GAME: '0',
        END_GAME: '2',
        PLAYERS: '3',
        PROMPT: 'a',
        DECISION_ACK: 'b',
        RESULT: 'c',
        WINNERS: 'd',
    }
}
Object.assign(en.ToServerCode, globalEn.ToServerCode);
Object.assign(en.ToClientCode, globalEn.ToClientCode);
export { en };