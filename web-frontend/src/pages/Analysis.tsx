import React from 'react';
import ConversationLayout from '../components/ConversationLayout';

const Analysis: React.FC = () => {
  return (
    <div className="analysis-page">
      <ConversationLayout conversationType="analysis" />
    </div>
  );
};

export default Analysis;
