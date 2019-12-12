import React from 'react';
import Container from './components/Container';
import "flexlayout-react/style/dark.css";
import './App.css';
import "normalize.css";
import "@blueprintjs/core/lib/css/blueprint.css";
import "@blueprintjs/icons/lib/css/blueprint-icons.css";
import { Button } from '@blueprintjs/core';
import Login from './components/Login';
import { Listing } from './serverapi/listing_pb';




const App: React.FC = () => {
  return (
    
      <Login ></Login>
    
  );
}

export default App;
