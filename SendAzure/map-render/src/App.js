import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';
import HeatMap from 'react-heatmap-grid';

const xLabels = new Array(50).fill(0).map((_, i) => "");
const yLabels = new Array(50).fill(0).map((_, i) => "");
const data = new Array(yLabels.length)
  .fill(0)
  .map(() => new Array(xLabels.length).fill(0).map(() => Math.floor(Math.random() * 100)));

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
        <div id="NW" className="quarter">
          <HeatMap
              xLabels={xLabels}
              yLabels={yLabels}
              data={data}
              height={4.5}
              xLabelWidth={1}
              background={"red"}
            />
        </div>

        <div id="NE" className="quarter">
          <HeatMap
              xLabels={xLabels}
              yLabels={yLabels}
              data={data}
              height={4.5}
              xLabelWidth={1}
              background={"blue"}
            />
        </div>

        <div id="SW" className="quarter">
          <HeatMap
              xLabels={xLabels}
              yLabels={yLabels}
              data={data}
              height={4.5}
              xLabelWidth={1}
              background={"green"}
            />
        </div>

        <div id="SE" className="quarter">
          <HeatMap
              xLabels={xLabels}
              yLabels={yLabels}
              data={data}
              height={4.5}
              xLabelWidth={1}
              background={"yellow"}
            />
        </div>

      </div>
    );
  }
}

export default App;
