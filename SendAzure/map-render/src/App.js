import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';
import HeatMap from './HeatMap';

const xLabels = new Array(50).fill(0).map((_, i) => "");
const yLabels = new Array(50).fill(0).map((_, i) => "");
let data1 = new Array(yLabels.length)
  .fill(100)
  .map(() => new Array(xLabels.length).fill(50));

let data2 = new Array(yLabels.length)
.fill(100)
.map(() => new Array(xLabels.length).fill(50));

let data3 = new Array(yLabels.length)
.fill(100)
.map(() => new Array(xLabels.length).fill(50));

let data4 = new Array(yLabels.length)
  .fill(100)
  .map(() => new Array(xLabels.length).fill(50));


let resultArray = [data1, data2, data3, data4];

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      allRobotsMap: {}
    };
  }

  mapPayloadToMap = () => {
    for (var key in this.state.allRobotsMap) {
      this.state.allRobotsMap[key].map((gridItem) => {
        if (gridItem.IsItFreeToRoam) {
          resultArray[Number(key)][gridItem.X+23][gridItem.Y+23] = 0
        } else {
          resultArray[Number(key)][gridItem.X+23][gridItem.Y+23] = 100
        }
      })
    }
  }


  componentDidMount(){
    setInterval(this.periodicFetchMap, 1000);
  }

  periodicFetchMap = () => {
    fetch("http://13.91.38.239:5000/getallmaps", {
      method: 'GET'
    })
    .then(res => res.json())
    .then(response => {
      console.log('this is the response')

      this.setState({allRobotsMap: response}, ()=> {
        this.mapPayloadToMap()
        // this.mapPoint()
      }
    )
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
              data={resultArray[0]}
              height={4.5}
              xLabelWidth={1}
              background={"red"}
            />
        </div>

        <div id="NE" className="quarter">
          <HeatMap
              xLabels={xLabels}
              yLabels={yLabels}
              data={resultArray[1]}
              height={4.5}
              xLabelWidth={1}
              background={"blue"}
            />
        </div>

        <div id="SW" className="quarter">
          <HeatMap
              xLabels={xLabels}
              yLabels={yLabels}
              data={resultArray[2]}
              height={4.5}
              xLabelWidth={1}
              background={"green"}
            />
        </div>

        <div id="SE" className="quarter">
          <HeatMap
              xLabels={xLabels}
              yLabels={yLabels}
              data={resultArray[3]}
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
