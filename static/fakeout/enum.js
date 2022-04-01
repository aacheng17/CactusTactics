import { globalEn } from '../globalEnum.js';

var en = {
    Phase: {
        PREGAME: '0',
        PLAY_PROMPT: '1',
        PLAY_GUESSES: '2',
    },
    ToServerCode: {
        END_GAME: '0',
        DECK_SELECTION: '1',
        SCORE_TO_WIN: '2',
        START_GAME: '3',
        RESPONSE: 'a',
        CHOICE: 'b',
        PROMPT_REQUEST: 'c',
    },
    ToClientCode: {
        START_GAME: '0',
        END_GAME: '1',
        PLAYERS: '2',
        DECK_OPTIONS: '4',
        DECK_SELECTION: '5',
        SCORE_TO_WIN: '6',
        PROMPT: 'a',
        CHOICE_RESPONSE: 'b',
        CHOICES: 'c',
        CHOICES_RESPONSE: 'd',
        RESULTS: 'e',
        WINNERS: 'f',
    }
}
Object.assign(en.ToServerCode, globalEn.ToServerCode);
Object.assign(en.ToClientCode, globalEn.ToClientCode);
export { en };