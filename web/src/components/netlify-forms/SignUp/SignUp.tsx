import React from 'react';

export default (
  <div className="uk-text-center margin-ends-48">
    <h3>SIGN UP FOR UPDATES</h3>
    <form
      name="contact"
      data-netlify
      className="uk-form-stacked">
      <p>
        <label className="uk-form-label">
          NAME
        </label>
        <input className="uk-input uk-width-medium" type="text" name="name" />
      </p>
      <p>
        <label className="uk-form-label">
          EMAIL
        </label>
        <input className="uk-input uk-width-medium" type="email" name="email" />
      </p>
      <p>
        <button type="submit" className="uk-button uk-button-default uk-width-medium">
          Send
        </button>
      </p>
    </form>
  </div>
)
