import React from 'react';

export default (
  <div className="uk-text-center margin-ends-48">
    <h3>SIGN UP FOR UPDATES</h3>
    <form
      name="contact"
      method="POST"
      className="uk-form-stacked"
    >
      <input type="hidden" name="form-name" value="contact" />
      <p>
        <label className="uk-form-label">
          NAME
        </label>
        <input type="text" name="name" className="uk-input uk-width-medium" />
      </p>
      <p>
        <label className="uk-form-label">
          EMAIL
        </label>
        <input type="email" name="email" className="uk-input uk-width-medium" />
      </p>
      <p>
        <button type="submit" className="uk-button uk-button-default uk-width-medium">
          Send
        </button>
      </p>
    </form>
  </div>
);
