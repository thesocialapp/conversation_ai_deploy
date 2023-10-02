// Chat.js
const Chat = ({ conversation }) => (
    <div>
      {conversation.map((entry, index) => (
        entry.me ? 
          <p key={index} className="user-message">Me: {entry.me}</p> : 
          <p key={index} className="bot-message">Bot: {entry.bot}</p>
      ))}
    </div>
  );
  
  export default Chat;
  