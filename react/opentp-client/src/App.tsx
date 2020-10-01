import "@blueprintjs/core/lib/css/blueprint.css";
import "@blueprintjs/icons/lib/css/blueprint-icons.css";
import "@blueprintjs/datetime/lib/css/blueprint-datetime.css";
import "flexlayout-react/style/dark.css";
import "normalize.css";
import React from 'react';
import './App.css';
import Login from './components/Login';
import log from 'loglevel';


const App: React.FC = () => {
  log.setLevel("INFO")
  return (

    <Login ></Login>

  );
}

export default App;
