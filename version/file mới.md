## Không có gì đâu

```python

@app.post("/chatbot/intent_recognition", tags=["Chatbot"])
async def intent_recognition(item: inputIntentRecognition):
    """API nhận diện intent theo format của hệ thống.

    ví dụ: 1
    ```
       {
        "question": "Xem chi tiết tk của tôi."
        }
    ```
    Returns:
     ```
     {
        "status": 1,
        "groupIntent": "",
        "intent": "API getAccountDetail",
        "options": [],
        "score": 0.8685,
        "question": "xem chi tiết tài khoản của tôi."
        }
     ```   
    """
    question = normalize_text(item.question) or ""
    question_normalized = vietnamese_normalizer.normalize_for_intent(question)
    itemSearchQuery=SearchQuery(query=question_normalized, top_k=config.DEFAULT_TOP_K)     
    
    KQ=search_single(query=itemSearchQuery)
    
    # results.append({
    #         'document': row['document'],
    #         'metadata': row['metadata'],
    #         'description': row['description'],
    #         'examples': row['examples_raw'],
    #         'keyword': row['keyword'],
    #         'score': round(float(row['score']), 4),
    #         'score_semantic': round(float(row['semantic_score']), 4),
    #         'score_max_example': round(float(row['max_example_score']), 4),
    #         'score_search_text': round(float(row['search_score']), 4),
    #         'score_bm25': round(float(row['bm25_normalized']), 4),
    #         'score_keyword_boost': round(float(row['keyword_boost']), 4),
    #         'best_matching_example': row.get('best_example', ''),
    #         'classified_intent': classified_intent
    #     })
    # print(KQ)
    
    return outputIntentRecognition(status=1, 
                                   groupIntent="", 
                                   intent=KQ[0]['document'], 
                                   options=[], 
                                   score=float(KQ[0]['score']), 
                                   question=question_normalized)

```


hết nội fung