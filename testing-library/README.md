# Mini Testing Library 🧪

Welcome to the Mini Testing Library, a simple yet powerful JavaScript testing library designed for developers who want to understand the inner workings of testing libraries.

## Example Usage 🚀
```ts
import { describe, expect, it } from './testing-library';
import { sum } from './sum';

describe('Basic addition', () => {
    it('adds 1 + 2 to equal 3', () => {
        expect(sum(1, 2)).toBe(3);
    });

    it('adds 2 + 5 to equal 7', () => {
        expect(sum(2, 5)).toBe(7);
    });
});
```

## Why building Testing Library? 🤔
> Building your own X gives you a chance to really understand how X works.

Step by step instructions on how to build this mini library, can be found in this [blog](https://mutasim.top/blog/build-testing-library).

Happy coding, and may all your test cases pass! 🥰
