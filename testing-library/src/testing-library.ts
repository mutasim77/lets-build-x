export function describe(suiteDescription: string, suiteFunction: () => void): void {
    console.log(`\nTest Suite: ${suiteDescription}`);

    try {
        suiteFunction();
    } catch (error) {
        throw new Error(`\nTest Suite "${suiteDescription}" failed`);
    }
}

export function it(description: string, test: () => void): void {
    console.log(`Test: ${description}`);
    try {
        test();
    } catch (error) {
        console.error(error);
        throw new Error('Test Run Failed');
    }
}

class Expectation<T> {
    private actual: T;

    constructor(actual: T) {
        this.actual = actual;
    }

    toBe(expected: unknown): void {
        if (this.actual === expected) {
            console.log('✓ Succeeded!');
        } else {
            throw new Error(`Fail - Expected ${this.actual} to be ${expected}`);
        }
    }

    toBeFalsy(): void {
        if (!this.actual) {
            console.log('✓ Succeeded!');
        } else {
            throw new Error(`Fail - Expected ${this.actual} to be falsy`);
        }
    }

    toBeTruthy(): void {
        if (this.actual) {
            console.log('✓ Succeeded!');
        } else {
            throw new Error(`Fail - Expected ${this.actual} to be truthy`);
        }
    }
}

export function expect<T>(actual: T): Expectation<T> {
    return new Expectation(actual);
}