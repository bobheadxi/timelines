import React, { ReactElement } from 'react';
import { AlertProps } from './index';

// TODO: add interactive elements
export default (props: AlertProps): ReactElement => {
  const { message } = props;
  return (
    <div className="uk-alert-warning" data-uk-alert>
      <p>{message}</p>
    </div>
  );
};
