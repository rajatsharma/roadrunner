class Person {
    constructor(name, age) {
        this.name = name;
        this.age = age;
    }

    greet(message) {
        return `${message}, I'm ${this.name}`;
    }

    birthday() {
        this.age += 1;
    }
}

function sum(x, y) {
    return x + y
}

const multiply = (x, y) => x * y;

const variable = 7;
