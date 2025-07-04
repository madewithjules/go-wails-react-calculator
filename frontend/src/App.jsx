import {useState, useEffect} from 'react';
import './App.css';
import {Calculate} from '../wailsjs/go/main/App'; // Ensure this path is correct

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
    ['x²', '√', 'π'],
  ];

  const getButtonClass = btn => {
    if (['/', '*', '-', '+', '=', '%', '√', 'x²'].includes(btn)) return 'button-operator';
    if (['C'].includes(btn)) return 'button-control';
    if (['(', ')', '.', 'π'].includes(btn)) return 'button-symbol';
    return 'button-number';
  };

  const handleButtonClick = async value => {
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
        setExpression(prev => (prev === '0' ? value : prev + '^2'));
        setDisplayValue(prev => (prev === '0' ? value : prev + '^2'));
        break;
      case '√':
        setExpression(prev => (prev === '0' ? 'sqrt(' : prev + 'sqrt('));
        setDisplayValue(prev => (prev === '0' ? 'sqrt(' : prev + 'sqrt('));
        break;
      case 'π':
        setExpression(prev => (prev === '0' ? 'π' : prev + 'π'));
        setDisplayValue(prev => (prev === '0' ? 'π' : prev + 'π'));
        break;
      case '.':
        const operatorChars = ['+', '-', '*', '/', '(', ')', '^', '%', 'sqrt('];
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
        setExpression(prev => (prev === '0' && value !== '(' && value !== ')' ? value : prev + value));
        setDisplayValue(prev => (prev === '0' && value !== '(' && value !== ')' ? value : prev + value));
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

  useEffect(() => {
    const handleKeyPress = (event) => {
      const {key} = event;
      let buttonValue = null;

      // Map keyboard keys to button values
      if (key >= '0' && key <= '9') {
        buttonValue = key;
      } else {
        switch (key) {
          case '+':
          case '-':
          case '*':
          case '/':
          case '%':
          case '(':
          case ')':
          case '.':
            buttonValue = key;
            break;
          case 'Enter':
          case '=':
            buttonValue = '=';
            break;
          case 'Backspace':
          case 'Delete':
            buttonValue = 'C';
            break;
          case 'Escape':
            buttonValue = 'C';
            break;
          // Add more mappings if necessary, e.g., for 'x²', '√', 'π'
          // These might require Shift or Alt modifiers, or different key choices
          // For simplicity, direct key mappings are shown here.
          // Example for 'x^2' (assuming 's' for square):
          // case 's':
          //   buttonValue = 'x²';
          //   break;
          // Example for 'sqrt' (assuming 'r' for root):
          // case 'r':
          //   buttonValue = '√';
          //   break;
          // Example for 'pi' (assuming 'p' for pi):
          // case 'p':
          //   buttonValue = 'π';
          //   break;
        }
      }

      if (buttonValue) {
        handleButtonClick(buttonValue);
      }
    };

    window.addEventListener('keydown', handleKeyPress);

    // Cleanup event listener on component unmount
    return () => {
      window.removeEventListener('keydown', handleKeyPress);
    };
  }, [handleButtonClick]); // Add handleButtonClick to dependency array if it's stable

  return (
    <div id="App">
      <div className="calculator">
        <div className="display">{displayValue}</div>
        <div className="buttons">
          {buttonLayout.flat().map((btn, i) => (
            <button key={i} className={`button ${getButtonClass(btn)}`} onClick={() => handleButtonClick(btn)}>
              {btn}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}

export default App;
