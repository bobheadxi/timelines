import React from 'react';
import DeckGL, { ScatterplotLayer } from 'deck.gl';

// Viewport settings
const viewState = {
  longitude: -74,
  latitude: 40.76,
  zoom: 13,
  maxZoom: 16,
  pitch: 50,
  bearing: 50
};

const MALE_COLOR = [0, 128, 255];
const FEMALE_COLOR = [255, 0, 128];

// DeckGL react component
class Burndown extends React.Component {
  render() {
    const layers = [
      new ScatterplotLayer({
        id: 'scatter-plot',
        data: 'https://raw.githubusercontent.com/uber-common/deck.gl-data/master/examples/scatterplot/manhattan.json',
        radiusScale: 10,
        radiusMinPixels: 0.5,
        getPosition: (d: any) => [d[0], d[1], 0],
        getColor: (d: any) => (d[2] === 1 ? MALE_COLOR : FEMALE_COLOR)
      })
    ];

    return (
      <DeckGL viewState={viewState} layers={layers} />
    );
  }
}

export default Burndown;
