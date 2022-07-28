from datetime import datetime
from scripts import get_public_attrs

# finish declaring models

class BaseModel:
    id:int= 0
    created_at:datetime = None
    updated_at:datetime = None
    

class Comic(BaseModel):
    name :str = ''
    description :str= ''
    status :str= ''
    number_of_episodes:int= ''
    cover_path:str= ''
    last_episode_time:datetime= None
    user_id :int= 0

class comic_episode(BaseModel):
    name:str = ''
    cover_path:str = ''
    episode_number:str= ''
    episode_path:str= ''
    user_id:int= 0
    comic_id:int= 0

class comic_comment(BaseModel):
    text = ''
    likes = 0
    dislikes = 0
    user_id = 0
    comic_id = 0

class comic_tag(BaseModel):
    name = ''