class Person {
  private name: any;
  private age: any;

  constructor(name: any, age: any) {
      this.name = name;
      this.age = age;
  }

  greet(message: any): any {
      return `${message}, I'm ${this.name}`;
  }

  birthday(): any {
      this.age += 1;
  }
}

function sum(x: any, y: any): any {
  return x + y
}

const multiply: any = (x: any, y: any): any => x * y;

const variable: any = 7;
