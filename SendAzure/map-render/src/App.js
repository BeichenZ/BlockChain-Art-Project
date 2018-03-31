import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      allRobotsMap: []
    };
  }

  periodicFetchMap = () => {
    fetch("http://localhost:8888/getmaps", {
      method: 'GET'
    })
    .then(res => res.json())
    .then(response => {
      console.log('this is the response')
      console.log(response)
      this.setState({allRobotsMap: response})
    })
    .catch(error => console.error('Error:', error))
  }

  render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h1 className="App-title">Welcome to React</h1>
        </header>
        <p className="App-intro">
          To get started, edit <code>src/App.js</code> and save to reload.
        </p>
      </div>
    );
  }
}

export default App;
