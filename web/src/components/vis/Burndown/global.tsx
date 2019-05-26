import React, { ReactElement } from 'react';
import DeckGL, { HexagonLayer } from 'deck.gl';

import { GlobalBurndown } from 'lib';

const viewState = {
  longitude: -122.4,
  latitude: 37.74,
  zoom: 11,
  maxZoom: 20,
  pitch: 30,
  bearing: 0,
};

export function globalBurndownElement(data: GlobalBurndown): ReactElement {
  const layers = [
    new HexagonLayer({
      id: 'scatter-plot',
      data: 'https://raw.githubusercontent.com/uber-common/deck.gl-data/master/website/sf-bike-parking.json',
      extruded: true,
      radius: 200,
      elevationScale: 4,
      getPosition: d => d.COORDINATES,
    }),
  ];

  return (
    <DeckGL viewState={viewState} layers={layers} />
  );
}
