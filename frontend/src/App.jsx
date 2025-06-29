import { useState } from 'react';
import './App.css';
import { Calculate } from '../wailsjs/go/main/App'; // Ensure this path is correct

function App() {
  const [expression, setExpression] = useState('0');
  const [displayValue, setDisplayValue] = useState('0');

  const buttonLayout = [
    ['C', '(', ')', '/'],
    ['7', '8', '9', '*'],
    ['4', '5', '6', '-'],
    ['1', '2', '3', '+'],
    ['0', '.', '%', '='],
    ['x²', '√', 'π']
  ];

  const getButtonClass = (btn) => {
    if (['/', '*', '-', '+', '=', '%', '√', 'x²'].includes(btn)) return 'button-operator';
    if (['C'].includes(btn)) return 'button-control';
    if (['(', ')', '.', 'π'].includes(btn)) return 'button-symbol';
    return 'button-number';
  };

  const handleButtonClick = async (value) => {
    switch (value) {
      case 'C':
        setExpression('0');
        setDisplayValue('0');
        break;
      case '=':
        try {
          // The backend Calculate function is currently a placeholder
          const result = await Calculate(expression);
          setDisplayValue(result);
          // Optional: setExpression(result) if you want the result to be the new starting point
          // For now, keeping expression as is until next calculation starts or 'C' is pressed
        } catch (error) {
          setDisplayValue('Error');
          console.error('Calculation error:', error);
        }
        break;
      case 'x²':
        setExpression(prev => prev === '0' ? value : prev + '^2');
        setDisplayValue(prev => prev === '0' ? value : prev + '^2');
        break;
      case '√':
        setExpression(prev => prev === '0' ? 'sqrt(' : prev + 'sqrt(');
        setDisplayValue(prev => prev === '0' ? 'sqrt(' : prev + 'sqrt(');
        break;
      case 'π':
        setExpression(prev => prev === '0' ? 'π' : prev + 'π');
        setDisplayValue(prev => prev === '0' ? 'π' : prev + 'π');
        break;
      case '.':
        // Avoid multiple dots in one number segment or starting with a dot if display is 0
        // This logic might need refinement for full correctness (e.g. "0.2.3" should not happen)
        if (expression === '0' || displayValue === '0') {
            setExpression('0.');
            setDisplayValue('0.');
        } else if (!expression.slice(expression.lastIndexOfAny(['+','-','*','/','(',')','^','%','sqrt(']) === -1 ? expression : expression.slice(expression.lastIndexOfAny(['+','-','*','/','(',')','^','%','sqrt('])+1)).includes('.')) {
            setExpression(prev => prev + '.');
            setDisplayValue(prev => prev + '.');
        }
        // A simpler dot logic for now: just append if expression is not '0'.
        // More complex logic would involve checking the current number being typed.
        // if (expression === '0') {
        //   setExpression('0.');
        //   setDisplayValue('0.');
        // } else if (!displayValue.includes('.')) { // Simple check on displayValue
        //   setExpression(prev => prev + '.');
        //   setDisplayValue(prev => prev + '.');
        // }
        break;
      case '+':
      case '-':
      case '*':
      case '/':
      case '%':
      case '(':
      case ')':
        // Append operator/symbol
        // If current expression is "0" (default), replace it. Else, append.
        setExpression(prev => (prev === '0' && value !== '(' && value !== ')') ? value : prev + value);
        setDisplayValue(prev => (prev === '0' && value !== '(' && value !== ')') ? value : prev + value);
        break;
      default: // Numbers
        if (expression === '0') {
          setExpression(value);
          setDisplayValue(value);
        } else {
          setExpression(prev => prev + value);
          setDisplayValue(prev => prev + value);
        }
        break;
    }
  };

// Standalone utility function
const findLastIndexOfAny = (str, chars) => {
  let lastIndex = -1;
  for (let i = 0; i < chars.length; i++) {
    const index = str.lastIndexOf(chars[i]);
    if (index > lastIndex) {
      lastIndex = index;
    }
  }
  return lastIndex;
};

function App() {
  const [expression, setExpression] = useState('0');
  const [displayValue, setDisplayValue] = useState('0');

  const buttonLayout = [
    ['C', '(', ')', '/'],
    ['7', '8', '9', '*'],
    ['4', '5', '6', '-'],
    ['1', '2', '3', '+'],
    ['0', '.', '%', '='],
    ['x²', '√', 'π']
  ];

  const getButtonClass = (btn) => {
    if (['/', '*', '-', '+', '=', '%', '√', 'x²'].includes(btn)) return 'button-operator';
    if (['C'].includes(btn)) return 'button-control';
    if (['(', ')', '.', 'π'].includes(btn)) return 'button-symbol';
    return 'button-number';
  };

  const handleButtonClick = async (value) => {
    switch (value) {
      case 'C':
        setExpression('0');
        setDisplayValue('0');
        break;
      case '=':
        try {
          const result = await Calculate(expression);
          setDisplayValue(result);
        } catch (error) {
          setDisplayValue('Error');
          console.error('Calculation error:', error);
        }
        break;
      case 'x²':
        setExpression(prev => prev === '0' ? value : prev + '^2');
        setDisplayValue(prev => prev === '0' ? value : prev + '^2');
        break;
      case '√':
        setExpression(prev => prev === '0' ? 'sqrt(' : prev + 'sqrt(');
        setDisplayValue(prev => prev === '0' ? 'sqrt(' : prev + 'sqrt(');
        break;
      case 'π':
        setExpression(prev => prev === '0' ? 'π' : prev + 'π');
        setDisplayValue(prev => prev === '0' ? 'π' : prev + 'π');
        break;
      case '.':
        const operatorChars = ['+','-','*','/','(',')','^','%','sqrt('];
        const lastOpIndex = findLastIndexOfAny(expression, operatorChars);
        const currentSegment = lastOpIndex === -1 ? expression : expression.slice(lastOpIndex + 1);

        if (expression === '0' || displayValue === '0') {
            setExpression('0.');
            setDisplayValue('0.');
        } else if (!currentSegment.includes('.')) {
            setExpression(prev => prev + '.');
            setDisplayValue(prev => prev + '.');
        }
        break;
      case '+':
      case '-':
      case '*':
      case '/':
      case '%':
      case '(':
      case ')':
        setExpression(prev => (prev === '0' && value !== '(' && value !== ')') ? value : prev + value);
        setDisplayValue(prev => (prev === '0' && value !== '(' && value !== ')') ? value : prev + value);
        break;
      default: // Numbers
        if (expression === '0') {
          setExpression(value);
          setDisplayValue(value);
        } else {
          setExpression(prev => prev + value);
          setDisplayValue(prev => prev + value);
        }
        break;
    }
  };

  return (
    <div id="App">
      <div className="calculator">
        <div className="display">{displayValue}</div>
        <div className="buttons">
          {buttonLayout.flat().map((btn, i) => (
            <button
              key={i}
              className={`button ${getButtonClass(btn)}`}
              onClick={() => handleButtonClick(btn)}
            >
              {btn}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}

export default App;
