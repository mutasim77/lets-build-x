import { describe, expect, it } from '../src/testing-library';
import { isEven, concatenateStrings } from './example';

describe('isEven function', () => {
    it('should be truthy for even numbers', () => {
        expect(isEven(6)).toBeTruthy();
        expect(isEven(8)).toBeTruthy();
    });

    it('should be falsy for odd numbers', () => {
        expect(isEven(1)).toBeFalsy();
    });
});

describe('Concatenate two strings', () => {
    it('will concatenate two strings', () => {
        const result = concatenateStrings('Hello', 'World');
        expect(result).toBe('HelloWorld');
    });

    it('should fail', () => {
        const result = concatenateStrings('2', '2');
        expect(result).toBe(4); //! '2' + '2' is '22' and 22 is not equal 4
    });
});